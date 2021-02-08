package findy

import (
	"time"

	"github.com/findy-network/findy-agent-api/grpc/agency"
	"github.com/findy-network/findy-agent-api/grpc/ops"
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

func (f *Agency) startHookOrWait() {
	const waitTime = 5
	for {
		err := f.listenAdminHook()
		if err == nil {
			break
		}
		glog.Warningf("listenAdminHook: cannot connect server, reconnecting after %d secs...", waitTime)
		time.Sleep(waitTime * time.Second)
	}
}

func (f *Agency) adminStatusLoop(ch chan *ops.AgencyStatus) {
	defer err2.Catch(func(err error) {
		glog.Errorf("Recovered error in psm hook routine: %s", err.Error())
		go f.adminStatusLoop(ch)
	})

	for {
		status, ok := <-ch
		if !ok {
			glog.Warningln("listenAdminHook: server lost, try reconnecting...")
			f.startHookOrWait()
			break
		}
		utils.LogMed().Infoln("received psm hook data for:", status.GetDID())

		protocolStatus := status.GetProtocolStatus()
		jobID := protocolStatus.State.ProtocolId.Id

		// TODO: pass also timestamps: when protocol was started/approved/sent/issued/verified etc.
		// revise this when we have "a real client" for the archive
		info := &model.ArchiveInfo{
			AgentID:       status.GetDID(),
			ConnectionID:  status.GetConnectionId(),
			JobID:         jobID,
			InitiatedByUs: protocolStatus.State.ProtocolId.Role == agency.Protocol_INITIATOR,
		}

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
}

func (f *Agency) listenAdminHook() (err error) {
	defer err2.Return(&err)

	glog.Info("Start listening to PSM events.")

	cmd := f.adminClient()
	// Error in registration is not notified here, instead all relevant info comes
	// in stream callback from now on
	ch, err := cmd.psmHook()
	err2.Check(err)

	go f.adminStatusLoop(ch)
	return nil
}
