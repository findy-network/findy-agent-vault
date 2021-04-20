package findy

import (
	"context"
	"io"

	"github.com/findy-network/findy-agent-api/grpc/agency"
	"github.com/findy-network/findy-agent-api/grpc/ops"
	"github.com/findy-network/findy-agent-vault/agency/model"
	"github.com/findy-network/findy-common-go/agency/client"
	"github.com/findy-network/findy-common-go/agency/client/async"
	"github.com/findy-network/findy-common-go/jwt"
	"github.com/golang/glog"
	"github.com/lainio/err2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/oauth"
)

type Client struct {
	*client.Conn
	ctx   context.Context
	cOpts []grpc.CallOption
}

func callOptions(jwtToken string) []grpc.CallOption {
	return []grpc.CallOption{
		grpc.PerRPCCredentials(
			oauth.NewOauthAccess(jwt.OauthToken(jwtToken)),
		),
	}
}

// Connection configuration for "sync" requests coming directly from web wallet
func (f *Agency) userSyncClient(a *model.Agent, connectionID string) *async.Pairwise {
	opts := callOptions(a.RawJWT)
	return async.NewPairwise(f.conn, connectionID, opts...)
}

// Connection configuration for "async" requests, done on behalf of the web wallet
func (f *Agency) userAsyncClient(a *model.Agent) *Client {
	opts := callOptions(jwt.BuildJWT(a.AgentID))
	return &Client{&f.conn, f.ctx, opts}
}

// Connection configuration for agency administrative client
func (f *Agency) adminClient() *Client {
	opts := callOptions(jwt.BuildJWT("findy-root"))
	return &Client{&f.conn, f.ctx, opts}
}

func (c *Client) release(id string, protocolType agency.Protocol_Type) (pid *agency.ProtocolID, err error) {
	protocolID := &agency.ProtocolID{
		Id:     id,
		TypeId: protocolType,
	}
	return c.Conn.DoRelease(c.ctx, protocolID, c.cOpts...)
}

func (c *Client) status(id string, protocolType agency.Protocol_Type) (pid *agency.ProtocolStatus, err error) {
	protocolID := &agency.ProtocolID{
		Id:     id,
		TypeId: protocolType,
	}
	return c.Conn.DoStatus(c.ctx, protocolID, c.cOpts...)
}

type AgentStatus struct {
	status *agency.AgentStatus
	err    error
}

func (c *Client) listen(id string) (ch chan *AgentStatus, err error) {
	clientID := &agency.ClientID{Id: id}
	defer err2.Return(&err)

	client := agency.NewAgentClient(c.ClientConn)
	statusCh := make(chan *AgentStatus)

	stream, err := client.Listen(c.ctx, clientID, c.cOpts...)
	err2.Check(err)
	glog.V(3).Infoln("successful start of listen id:", clientID.Id)

	go func() {
		defer err2.CatchTrace(func(err error) {
			glog.V(1).Infoln("WARNING: error when reading response:", err)
			close(statusCh)
		})
		for {
			status, err := stream.Recv()
			if err == io.EOF {
				glog.V(3).Infoln("status stream end")
				close(statusCh)
				break
			}
			statusCh <- &AgentStatus{status, err}
			err2.Check(err)
		}
	}()
	return statusCh, nil
}

func (c *Client) psmHook() (ch chan *ops.AgencyStatus, err error) {
	return c.Conn.PSMHook(c.ctx, c.cOpts...)
}
