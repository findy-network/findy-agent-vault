package findy

import (
	"github.com/findy-network/findy-agent-api/grpc/agency"
	"github.com/findy-network/findy-agent-vault/agency/model"
	"github.com/findy-network/findy-agent-vault/utils"
	"github.com/golang/glog"
	"github.com/lainio/err2"
)

func (f *Agency) archive(info *model.ArchiveInfo, status *agency.ProtocolStatus) {
	switch status.State.ProtocolId.TypeId {
	case agency.Protocol_CONNECT:
		connection := statusToConnection(status)
		f.archiver.ArchiveConnection(info, connection)
	case agency.Protocol_ISSUE:
		credential := statusToCredential(status)
		f.archiver.ArchiveCredential(info, credential)
	case agency.Protocol_PROOF:
		proof := statusToProof(status)
		f.archiver.ArchiveProof(info, proof)
	case agency.Protocol_BASIC_MESSAGE:
		message := statusToMessage(status)
		f.archiver.ArchiveMessage(info, message)
	default:
		utils.LogHigh().Infof(
			"Received unknown protocol type %s",
			status.State.ProtocolId.TypeId.String(),
		)
	}
}

func (f *Agency) listenAdminHook() (err error) {
	defer err2.Return(&err)

	// TODO: cancellation, reconnect
	glog.Info("Start listening to PSM events.")

	cmd := f.adminClient()
	// Error in registration is not notified here, instead all relevant info comes
	// in stream callback from now on
	ch, err := cmd.psmHook()
	err2.Check(err)

	go func() {
		defer err2.Catch(func(err error) {
			glog.Errorf("Recovered error in psm hook routine: %s", err.Error())
			// TODO: reconnect?
		})

		for {
			status, ok := <-ch
			if !ok {
				glog.Warningln("closed from server")
				break
			}
			utils.LogMed().Infoln("received psm hook data for:", status.GetDID())

			protocolStatus := status.GetProtocolStatus()
			info := &model.ArchiveInfo{AgentID: status.GetDID(), ConnectionID: status.GetConnectionId()}

			// archive currently only successful protocol results
			if protocolStatus.State.State == agency.ProtocolState_OK {
				f.archive(info, protocolStatus)
			} else {
				utils.LogLow().Infof(
					"Skipping archiving for protocol run %s in state %s",
					protocolStatus.State.ProtocolId.TypeId,
					protocolStatus.State.State,
				)
			}
		}
	}()
	return nil
}
