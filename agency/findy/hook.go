package findy

import (
	"github.com/findy-network/findy-agent-api/grpc/agency"
	"github.com/findy-network/findy-agent-vault/agency/model"
	"github.com/findy-network/findy-agent-vault/utils"
	"github.com/findy-network/findy-grpc/agency/client"
	auth "github.com/findy-network/findy-grpc/jwt"
	"github.com/findy-network/findy-grpc/rpc"
	"github.com/golang/glog"
	"github.com/lainio/err2"
	"google.golang.org/grpc"
)

func (f *Agency) adminClient(user string) (conn *grpc.ClientConn, err error) {
	defer err2.Return(&err)

	utils.LogLow().Infoln("client with user:", user)

	cfg := client.BuildClientConnBase(f.tlsPath, f.agencyHost, f.agencyPort, f.options)
	token := auth.BuildJWT(user)
	cfg.JWT = token

	conn, err = rpc.ClientConn(*cfg)
	err2.Check(err)
	return
}

func (f *Agency) adminListenClient() client.Conn {
	config := client.BuildClientConnBase(f.tlsPath, f.agencyHost, f.agencyPort, f.options)
	return client.TryOpen("findy-root", config)
}

func (f *Agency) listenAdminHook() (err error) {
	defer err2.Return(&err)

	// TODO: cancellation, reconnect
	glog.Info("Start listening to PSM events.")

	conn := f.adminListenClient()
	// Error in registration is not notified here, instead all relevant info comes
	// in stream callback from now on
	ch, err := conn.PSMHook(f.ctx)
	err2.Check(err)

	go func() {
		defer err2.Catch(func(err error) {
			glog.Errorf("Recovered error in psm hook routine: %s", err.Error())
			// TODO: reconnect?
		})

		// TODO: fail job if error happens?
		for {
			status, ok := <-ch
			if !ok {
				glog.Warningln("closed from server")
				conn.Close()
				break
			}
			utils.LogMed().Infoln("received psm hook data for:", status.GetDID())

			protocolStatus := status.GetProtocolStatus()
			info := &model.ArchiveInfo{AgentID: status.GetDID(), ConnectionID: ""}

			// archive currently only successful protocol runs
			if protocolStatus.State.State == agency.ProtocolState_OK {
				switch status.ProtocolStatus.State.ProtocolId.TypeId {
				case agency.Protocol_CONNECT:
					connection := statusToConnection(protocolStatus)
					f.archiver.ArchiveConnection(info, connection)
				case agency.Protocol_ISSUE:
					credential := statusToCredential(protocolStatus)
					f.archiver.ArchiveCredential(info, credential)
				case agency.Protocol_PROOF:
					proof := statusToProof(protocolStatus)
					f.archiver.ArchiveProof(info, proof)
				case agency.Protocol_BASIC_MESSAGE:
					message := statusToMessage(protocolStatus)
					f.archiver.ArchiveMessage(info, message)
				}
			} else {
				utils.LogLow().Infof(
					"Skipping archiving for protocol run in state %s",
					protocolStatus.State.State,
				)
			}

		}
	}()

	/*
		conn, err := f.adminClient("findy-root")

		err2.Check(err)
		opsClient := ops.NewAgencyClient(conn)

		statusCh := make(chan *agency.ProtocolStatus)

		// Error in registration is not notified here, instead all relevant info comes
		// in stream callback from now on
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
				status, ok := <-statusCh
				if !ok {
					glog.Warning("closed from server")
					break
				}
				// TODO: get agent ID from agency status
				// store data
			}
		}()
	*/
	return nil
}
