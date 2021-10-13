package findy

import (
	"time"

	"github.com/findy-network/findy-agent-vault/agency/model"
	"github.com/findy-network/findy-agent-vault/utils"
	agency "github.com/findy-network/findy-common-go/grpc/agency/v1"
	"github.com/golang/glog"
	"github.com/lainio/err2"
	"google.golang.org/grpc/codes"
	grpcStatus "google.golang.org/grpc/status"
)

type counter struct {
	count    int
	lastCode codes.Code
}

func (c *counter) reset() {
	c.count = 0
	c.lastCode = codes.Unknown
}

func (f *Agency) getStatus(a *model.Agent, notification *agency.Notification) (status *agency.ProtocolStatus, ok bool) {
	cmd := f.userAsyncClient(a)

	status, err := cmd.status(notification.ProtocolID, notification.ProtocolType)

	if err != nil {
		glog.Errorf("Unable to fetch protocol status for %s (%s)", notification.ProtocolID, err.Error())
		return
	}

	if status == nil {
		glog.Errorf("Received invalid protocol status for %s", notification.ProtocolID)
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
	case agency.Protocol_ISSUE_CREDENTIAL:
		err2.Check(f.vault.UpdateCredential(
			job,
			&model.CredentialUpdate{
				FailedMs: &now,
			},
		))
	case agency.Protocol_PRESENT_PROOF:
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
	case agency.Protocol_DIDEXCHANGE:
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

	case agency.Protocol_ISSUE_CREDENTIAL:
		err2.Check(f.vault.UpdateCredential(
			job,
			&model.CredentialUpdate{
				IssuedMs: &now,
			},
		))
	case agency.Protocol_PRESENT_PROOF:
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
			f.releaseCompleted(a, status.State.ProtocolID.ID, status.State.ProtocolID.TypeID)
		}
	case agency.ProtocolState_OK:
		if f.handleProtocolSuccess(job, notification, status) == nil {
			f.releaseCompleted(a, status.State.ProtocolID.ID, status.State.ProtocolID.TypeID)
		}
	default:
		utils.LogLow().Infof(
			"Received status update %s: %s",
			status.State.ProtocolID.GetTypeID().String(),
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
	case agency.Protocol_ISSUE_CREDENTIAL:
		credential := statusToCredential(status)
		if credential == nil {
			glog.Errorf("Received invalid credential issue object for %s", job.JobID)
			return
		}
		// TODO: what if we are issuer?
		_ = f.vault.AddCredential(job, credential)

	case agency.Protocol_PRESENT_PROOF:
		proof := statusToProof(status)
		if proof == nil {
			glog.Errorf("Received invalid proof object for %s", job.JobID)
			return
		}
		// TODO: what if we are verifier?
		_ = f.vault.AddProof(job, proof)

	case agency.Protocol_NONE:
	case agency.Protocol_TRUST_PING:
	case agency.Protocol_DIDEXCHANGE:
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
	switch notification.TypeID {
	case agency.Notification_PROTOCOL_PAUSED:
		f.handleAction(job, notification, status)
	case agency.Notification_STATUS_UPDATE:
		f.handleStatus(a, job, notification, status)
	case agency.Notification_NONE:
	case agency.Notification_KEEPALIVE:
		// TODO?
	}
}

func (f *Agency) waitAndRetryListening(a *model.Agent, err error, retryCounter counter) counter {
	const maxCount = 5
	const waitTime = 5
	count := retryCounter.count

	utils.LogLow().Infoln("Listen and wait", count)

	errCode := codes.Unknown
	if e, ok := grpcStatus.FromError(err); ok {
		errCode = e.Code()
	}

	if count < maxCount || errCode == codes.Unavailable { // try to connect infinitely if service is coming up
		glog.Warningln("listenAgent: channel closed, try reconnecting...", count)
		if errCode == retryCounter.lastCode {
			count++
		} else {
			count = 0
		}
		for {
			newWaitTime := count * waitTime
			glog.Warningf("listenAgent: waiting, reconnecting after %d secs...", newWaitTime)
			time.Sleep(time.Duration(newWaitTime) * time.Second)

			err := f.listenAgentWithRetry(a, counter{count, errCode})
			if err == nil {
				utils.LogLow().Infoln("Agent listening retry succeeded.")
				break
			}
			glog.Warningf("listenAgent: cannot connect server, try again...")
		}
	} else {
		glog.Errorf("Retry count exceeded with code %v, aborting listening for %v.", errCode, a.AgentID)
	}

	return counter{count, errCode}
}

func (f *Agency) agentStatusLoop(a *model.Agent, ch chan *AgentStatus, retryCounter counter) {
	defer err2.Catch(func(err error) {
		glog.Errorf("Recovered error in agent listener routine: %s, continue listening...", err.Error())

		go f.agentStatusLoop(a, ch, counter{})
	})

	utils.LogLow().Infoln("Start agentStatusLoop for", a.AgentID)

	for {
		chRes, ok := <-ch
		var chErr error
		if chRes != nil {
			chErr = chRes.err
		}
		if !ok || chErr != nil {
			f.waitAndRetryListening(a, chErr, retryCounter)
			break
		}

		status := chRes.status

		if status.Notification == nil {
			glog.Warningf("Received status with no notification: %+v", status)
			continue
		}

		// successful round -> reset retry counter
		retryCounter.reset()

		if status.Notification.TypeID == agency.Notification_KEEPALIVE {
			utils.LogTrace().Infof("Keepalive for agent %s", a.TenantID)
			continue
		}

		utils.LogMed().Infoln("received notification:",
			status.Notification.TypeID,
			status.Notification.Role,
			status.Notification.ProtocolID)

		job := &model.JobInfo{
			TenantID:     a.TenantID,
			JobID:        status.Notification.ProtocolID,
			ConnectionID: status.Notification.ConnectionID,
		}

		protocolStatus, ok := f.getStatus(a, status.Notification)
		if !ok {
			continue
		}

		f.handleNotification(a, job, status.Notification, protocolStatus)
	}
}

func (f *Agency) listenAgent(a *model.Agent) (err error) {
	return f.listenAgentWithRetry(a, counter{})
}

func (f *Agency) listenAgentWithRetry(a *model.Agent, retryCounter counter) (err error) {
	defer err2.Return(&err)

	utils.LogLow().Infoln("Listen agent with retry count", retryCounter.count)

	cmd := f.userAsyncClient(a)

	// Error in registration is not notified here, instead all relevant info comes
	// in stream callback from now on
	ch, err := cmd.listen(a.TenantID)
	err2.Check(err)

	go f.agentStatusLoop(a, ch, retryCounter)

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
