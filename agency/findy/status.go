package findy

import (
	"time"

	"github.com/findy-network/findy-agent-api/grpc/agency"
	"github.com/findy-network/findy-agent-vault/agency/model"
	"github.com/findy-network/findy-agent-vault/utils"
	"github.com/golang/glog"
	"github.com/lainio/err2"
)

func (f *Agency) getStatus(a *model.Agent, notification *agency.Notification) (status *agency.ProtocolStatus, ok bool) {
	cmd := f.userAsyncClient(a)

	status, err := cmd.status(notification.ProtocolId, notification.ProtocolType)

	if err != nil {
		glog.Errorf("Unable to fetch protocol status for %s (%s)", notification.ProtocolId, err.Error())
		return
	}

	if status == nil {
		glog.Errorf("Received invalid protocol status for %s", notification.ProtocolId)
		return
	}

	ok = true
	return
}

func (f *Agency) handleProtocolFailure(
	job *model.JobInfo,
	notification *agency.Notification,
) (err error) {
	err2.Return(&err)

	// TODO: failure reason
	utils.LogHigh().Infof("Job %s (%s) failed", job.JobID, notification.ProtocolType.String())

	now := f.currentTimeMs()
	switch notification.ProtocolType {
	case agency.Protocol_ISSUE:
		err2.Check(f.vault.UpdateCredential(
			job,
			&model.CredentialUpdate{
				FailedMs: &now,
			},
		))
	case agency.Protocol_PROOF:
		err2.Check(f.vault.UpdateProof(
			job,
			&model.ProofUpdate{
				FailedMs: &now,
			},
		))
	default:
		err2.Check(f.vault.FailJob(job))
	}
	return
}

func (f *Agency) handleProtocolSuccess(
	job *model.JobInfo,
	notification *agency.Notification,
	status *agency.ProtocolStatus,
) (err error) {
	err2.Return(&err)

	utils.LogLow().Infof("Job %s (%s) success", job.JobID, notification.ProtocolType.String())

	now := f.currentTimeMs()
	switch notification.ProtocolType {
	case agency.Protocol_CONNECT:
		connection := statusToConnection(status)
		if connection == nil {
			glog.Errorf("Received invalid connection object for %s", job.JobID)
			return
		}

		err2.Check(f.vault.AddConnection(job, connection))

	case agency.Protocol_BASIC_MESSAGE:
		message := statusToMessage(status)
		if message == nil {
			glog.Errorf("Received invalid message object for %s", job.JobID)
			return
		}

		// TODO: delivered?
		err2.Check(f.vault.AddMessage(job, message))

	case agency.Protocol_ISSUE:
		err2.Check(f.vault.UpdateCredential(
			job,
			&model.CredentialUpdate{
				IssuedMs: &now,
			},
		))
	case agency.Protocol_PROOF:
		err2.Check(f.vault.UpdateProof(
			job,
			&model.ProofUpdate{
				VerifiedMs: &now,
			},
		))
	case agency.Protocol_NONE:
	case agency.Protocol_TRUST_PING:
	}

	return nil
}

func (f *Agency) handleStatus(
	a *model.Agent,
	job *model.JobInfo,
	notification *agency.Notification,
	status *agency.ProtocolStatus,
) {
	switch status.State.State {
	case agency.ProtocolState_ERR:
		if f.handleProtocolFailure(job, notification) == nil {
			f.releaseCompleted(a, status.State.ProtocolId.Id, status.State.ProtocolId.TypeId)
		}
	case agency.ProtocolState_OK:
		if f.handleProtocolSuccess(job, notification, status) == nil {
			f.releaseCompleted(a, status.State.ProtocolId.Id, status.State.ProtocolId.TypeId)
		}
	default:
		utils.LogLow().Infof(
			"Received status update %s: %s",
			status.State.ProtocolId.GetTypeId().String(),
			status.State.GetState().String(),
		)
	}
}

func (f *Agency) handleAction(
	job *model.JobInfo,
	notification *agency.Notification,
	status *agency.ProtocolStatus,
) {
	switch notification.ProtocolType {
	case agency.Protocol_ISSUE:
		credential := statusToCredential(status)
		if credential == nil {
			glog.Errorf("Received invalid credential issue object for %s", job.JobID)
			return
		}
		// TODO: what if we are issuer?
		_ = f.vault.AddCredential(job, credential)

	case agency.Protocol_PROOF:
		proof := statusToProof(status)
		if proof == nil {
			glog.Errorf("Received invalid proof object for %s", job.JobID)
			return
		}
		// TODO: what if we are verifier?
		_ = f.vault.AddProof(job, proof)

	case agency.Protocol_NONE:
	case agency.Protocol_TRUST_PING:
	case agency.Protocol_CONNECT:
	case agency.Protocol_BASIC_MESSAGE:
		// N/A
		glog.Errorf("Should not handle action for protocol %s", notification.ProtocolType)
	}
}

func (f *Agency) handleNotification(
	a *model.Agent,
	job *model.JobInfo,
	notification *agency.Notification,
	status *agency.ProtocolStatus,
) {
	switch notification.TypeId {
	case agency.Notification_ACTION_NEEDED:
		f.handleAction(job, notification, status)
	case agency.Notification_STATUS_UPDATE:
		f.handleStatus(a, job, notification, status)
	case agency.Notification_ANSWER_NEEDED_PING:
	case agency.Notification_ANSWER_NEEDED_ISSUE_PROPOSE:
	case agency.Notification_ANSWER_NEEDED_PROOF_PROPOSE:
	case agency.Notification_ANSWER_NEEDED_PROOF_VERIFY:
		// TODO?
	}
}

func (f *Agency) startListeningOrWait(a *model.Agent, retryCount int) int {
	const maxCount = 5
	const waitTime = 5

	if retryCount < maxCount {
		glog.Warningln("listenAgent: channel closed, try reconnecting...", retryCount)
		retryCount++
		for {
			err := f.listenAgentWithRetry(a, retryCount)
			if err == nil {
				break
			}
			glog.Warningf("listenAgent: cannot connect server, reconnecting after %d secs...", waitTime)
			time.Sleep(waitTime * time.Second)
		}
	}

	return retryCount
}

func (f *Agency) agentStatusLoop(a *model.Agent, ch chan *agency.AgentStatus, retryCount int) {
	defer err2.Catch(func(err error) {
		glog.Errorf("Recovered error in agent listener routine: %s, continue listening...", err.Error())

		go f.agentStatusLoop(a, ch, 0)
	})

	for {
		status, ok := <-ch
		if !ok {
			retryCount = f.startListeningOrWait(a, retryCount)
			break
		}

		if status.Notification == nil {
			glog.Warningf("Received status with no notification: %+v", status)
			continue
		}

		// successful round -> reset retry counter
		retryCount = 0

		utils.LogMed().Infoln("received notification:",
			status.Notification.TypeId,
			status.Notification.Role,
			status.Notification.ProtocolId)

		job := &model.JobInfo{
			TenantID:     a.TenantID,
			JobID:        status.Notification.ProtocolId,
			ConnectionID: status.Notification.ConnectionId,
		}

		protocolStatus, ok := f.getStatus(a, status.Notification)
		if !ok {
			continue
		}

		f.handleNotification(a, job, status.Notification, protocolStatus)
	}
}

func (f *Agency) listenAgent(a *model.Agent) (err error) {
	return f.listenAgentWithRetry(a, 0)
}

func (f *Agency) listenAgentWithRetry(a *model.Agent, retryCount int) (err error) {
	defer err2.Return(&err)

	cmd := f.userAsyncClient(a)

	// Error in registration is not notified here, instead all relevant info comes
	// in stream callback from now on
	ch, err := cmd.listen(a.TenantID)
	err2.Check(err)

	go f.agentStatusLoop(a, ch, retryCount)

	return err
}

func (f *Agency) releaseCompleted(a *model.Agent, protocolID string, protocolType agency.Protocol_Type) {
	defer err2.Catch(func(err error) {
		glog.Errorf("Failure when releasing protocol: %s", err.Error())
	})

	cmd := f.userAsyncClient(a)
	_, err := cmd.release(protocolID, protocolType)
	err2.Check(err)
}
