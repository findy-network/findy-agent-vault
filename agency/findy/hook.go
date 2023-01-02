package findy

import (
	"time"

	"github.com/findy-network/findy-agent-vault/agency/model"
	"github.com/findy-network/findy-agent-vault/utils"
	agency "github.com/findy-network/findy-common-go/grpc/agency/v1"
	ops "github.com/findy-network/findy-common-go/grpc/ops/v1"
	"github.com/golang/glog"
	"github.com/lainio/err2"
	"github.com/lainio/err2/try"
)

const waitTime = 5

func (f *Agency) archive(info *model.ArchiveInfo, status *agency.ProtocolStatus) {
	switch status.State.ProtocolID.TypeID {
	case agency.Protocol_DIDEXCHANGE:
		connection := statusToConnection(status)
		f.archiver.ArchiveConnection(info, connection)
	case agency.Protocol_ISSUE_CREDENTIAL:
		credential := statusToCredential(status)
		f.archiver.ArchiveCredential(info, credential)
	case agency.Protocol_PRESENT_PROOF:
		proof := statusToProof(status)
		f.archiver.ArchiveProof(info, proof)
	case agency.Protocol_BASIC_MESSAGE:
		message := statusToMessage(status)
		f.archiver.ArchiveMessage(info, message)
	default:
		utils.LogHigh().Infof(
			"Received unknown protocol type %s",
			status.State.ProtocolID.TypeID.String(),
		)
	}
}

func (f *Agency) startHookOrWait() {
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
			time.Sleep(waitTime * time.Second)
			f.startHookOrWait()
			break
		}
		utils.LogMed().Infoln("received psm hook data for:", status.GetDID())

		protocolStatus := status.GetProtocolStatus()
		jobID := protocolStatus.State.ProtocolID.ID

		// TODO: pass also timestamps: when protocol was started/approved/sent/issued/verified etc.
		// revise this when we have "a real client" for the archive
		info := &model.ArchiveInfo{
			AgentID:       status.GetDID(),
			ConnectionID:  status.GetConnectionID(),
			JobID:         jobID,
			InitiatedByUs: protocolStatus.State.ProtocolID.Role == agency.Protocol_INITIATOR,
		}

		// archive currently only successful protocol results
		if protocolStatus.State.State == agency.ProtocolState_OK {
			f.archive(info, protocolStatus)
		} else {
			utils.LogLow().Infof(
				"Skipping archiving for protocol run %s in state %s",
				protocolStatus.State.ProtocolID.TypeID,
				protocolStatus.State.State,
			)
		}
	}
}

func (f *Agency) listenAdminHook() (err error) {
	defer err2.Handle(&err)

	glog.Info("Start listening to PSM events.")

	cmd := f.adminClient()
	// Error in registration is not notified here, instead all relevant info comes
	// in stream callback from now on
	ch := try.To1(cmd.psmHook())

	go f.adminStatusLoop(ch)
	return nil
}
