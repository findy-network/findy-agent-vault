package findy

import (
	"context"

	"github.com/findy-network/findy-agent-api/grpc/agency"
	"github.com/findy-network/findy-agent-api/grpc/ops"
	"github.com/findy-network/findy-agent-vault/agency/model"
	"github.com/findy-network/findy-grpc/agency/client"
	"github.com/findy-network/findy-grpc/agency/client/async"
	"github.com/findy-network/findy-grpc/jwt"
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

// Connection configuration for agency administrative clietn
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

func (c *Client) listen(id string) (ch chan *agency.AgentStatus, err error) {
	clientID := &agency.ClientID{Id: id}
	return c.Conn.Listen(c.ctx, clientID, c.cOpts...)
}

func (c *Client) psmHook() (ch chan *ops.AgencyStatus, err error) {
	return c.Conn.PSMHook(c.ctx, c.cOpts...)
}
