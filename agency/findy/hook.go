package findy

import (
	"fmt"
	"io"

	"github.com/findy-network/findy-agent-api/grpc/agency"
	"github.com/findy-network/findy-agent-api/grpc/ops"
	"github.com/findy-network/findy-grpc/agency/client"
	auth "github.com/findy-network/findy-grpc/jwt"
	"github.com/findy-network/findy-grpc/rpc"
	"github.com/findy-network/findy-grpc/utils"
	"github.com/golang/glog"
	"github.com/lainio/err2"
	"google.golang.org/grpc"
)

func adminClient(user string) (conn *grpc.ClientConn, err error) {
	defer err2.Return(&err)

	glog.V(5).Infoln("client with user:", user)

	cfg := client.BuildClientConnBase("", agencyHost, agencyPort, nil)
	token := auth.BuildJWT(user)
	cfg.JWT = token

	conn, err = rpc.ClientConn(*cfg)
	err2.Check(err)
	return
}

func (f *Agency) listenAdminHook() (err error) {
	glog.Info("Start listening to PSM events.")

	conn, err := adminClient("findy-root")
	err2.Check(err)
	opsClient := ops.NewAgencyClient(conn)

	statusCh := make(chan *agency.ProtocolStatus)

	stream, err := opsClient.PSMHook(f.ctx, &ops.DataHook{Id: utils.UUID()})
	err2.Check(err)
	glog.V(3).Infoln("successful start of listen PSM hook id:")

	go func() {
		defer err2.CatchTrace(func(err error) {
			glog.V(1).Infoln("WARNING: error when reading response:", err)
			close(statusCh)
			conn.Close()
			// TODO: reconnect logic
		})
		for {
			status, err := stream.Recv()
			if err == io.EOF {
				glog.V(3).Infoln("status stream end")
				close(statusCh)
				conn.Close()
				break
			}
			err2.Check(err)
			fmt.Println(status)
			statusCh <- status.ProtocolStatus
		}
	}()

	go func() {
		for {
			status, ok := <-statusCh
			if !ok {
				glog.V(2).Infoln("closed from server")
				break
			}
			switch status.State.ProtocolId.TypeId {
			case agency.Protocol_CONNECT:
				// connection := status.GetConnection()
				// TODO: get agent ID from agency status
				// store data
			default:
				fmt.Println("Skip:", status.State.ProtocolId.TypeId.String())
			}

		}
	}()

	return nil
}
