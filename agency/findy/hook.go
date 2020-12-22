// +build findy_grpc

package findy

import (
	"io"

	"github.com/findy-network/findy-agent-api/grpc/agency"
	"github.com/findy-network/findy-agent-api/grpc/ops"
	"github.com/findy-network/findy-agent-vault/utils"
	"github.com/findy-network/findy-grpc/agency/client"
	auth "github.com/findy-network/findy-grpc/jwt"
	"github.com/findy-network/findy-grpc/rpc"
	"github.com/golang/glog"
	"github.com/google/uuid"
	"github.com/lainio/err2"
	"google.golang.org/grpc"
)

func adminClient(user string) (conn *grpc.ClientConn, err error) {
	defer err2.Return(&err)

	utils.LogLow().Infoln("client with user:", user)

	cfg := client.BuildClientConnBase("", agencyHost, agencyPort, nil)
	token := auth.BuildJWT(user)
	cfg.JWT = token

	conn, err = rpc.ClientConn(*cfg)
	err2.Check(err)
	return
}

func (f *Agency) listenAdminHook() (err error) {
	defer err2.Return(&err)

	// TODO: cancellation, reconnect
	glog.Info("Start listening to PSM events.")

	conn, err := adminClient("findy-root")
	err2.Check(err)
	opsClient := ops.NewAgencyClient(conn)

	statusCh := make(chan *agency.ProtocolStatus)

	stream, err := opsClient.PSMHook(f.ctx, &ops.DataHook{Id: uuid.New().String()})
	err2.Check(err)
	utils.LogMed().Infoln("successful start of listen PSM hook id:")

	go func() {
		defer err2.CatchTrace(func(err error) {
			glog.Warningln("error when reading response:", err)
			close(statusCh)
			conn.Close()
			// TODO: reconnect logic
		})
		for {
			status, err := stream.Recv()
			if err == io.EOF {
				glog.Warningln("status stream end")
				close(statusCh)
				conn.Close()
				break
			}
			err2.Check(err)
			statusCh <- status.ProtocolStatus
		}
	}()

	go func() {
		for {
			_, ok := <-statusCh
			if !ok {
				glog.Warning("closed from server")
				break
			}
			// TODO: get agent ID from agency status
			// store data
		}
	}()

	return nil
}
