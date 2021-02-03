package findy

import (
	"context"

	"github.com/findy-network/findy-agent-api/grpc/agency"
	"github.com/findy-network/findy-agent-vault/agency/model"
	"github.com/findy-network/findy-agent-vault/utils"
	"github.com/findy-network/findy-grpc/agency/client"
	"github.com/golang/glog"
	"github.com/lainio/err2"
)

func (f *Agency) userListenClient(a *model.Agent) client.Conn {
	config := client.BuildClientConnBase(f.tlsPath, f.agencyHost, f.agencyPort, f.options)
	return client.TryOpen(a.AgentID, config)
}

func (f *Agency) getStatus(conn client.Conn, notification *agency.Notification) (status *agency.ProtocolStatus, ok bool) {
	var err error

	ctx := context.Background()
	didComm := agency.NewDIDCommClient(conn)
	status, err = didComm.Status(ctx, &agency.ProtocolID{
		TypeId:           notification.ProtocolType,
		Id:               notification.ProtocolId,
		NotificationTime: notification.Timestamp,
	})

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

	now := utils.CurrentTimeMs()
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

	now := utils.CurrentTimeMs()
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
	conn client.Conn,
	job *model.JobInfo,
	notification *agency.Notification,
	status *agency.ProtocolStatus,
) {
	switch status.State.State {
	case agency.ProtocolState_ERR:
		if f.handleProtocolFailure(job, notification) == nil {
			f.releaseCompleted(conn, job, notification, status)
		}
	case agency.ProtocolState_OK:
		if f.handleProtocolSuccess(job, notification, status) == nil {
			f.releaseCompleted(conn, job, notification, status)
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

func (f *Agency) listenAgent(a *model.Agent) (err error) {
	defer err2.Return(&err)
	// TODO: cancellation, reconnect

	conn := f.userListenClient(a)

	// Error in registration is not notified here, instead all relevant info comes
	// in stream callback from now on
	ch, err := conn.Listen(f.ctx, &agency.ClientID{Id: a.TenantID})
	err2.Check(err)

	go func() {
		defer err2.Catch(func(err error) {
			glog.Errorf("Recovered error in listener routine: %s", err.Error())
			// TODO: reconnect?
		})

		for {
			status, ok := <-ch
			if !ok {
				glog.Warningln("closed from server")
				conn.Close()
				break
			}
			utils.LogMed().Infoln("received notification:",
				status.Notification.TypeId,
				status.Notification.Role,
				status.Notification.ProtocolId)

			job := &model.JobInfo{
				TenantID:     a.TenantID,
				JobID:        status.Notification.ProtocolId,
				ConnectionID: status.Notification.ConnectionId,
			}

			protocolStatus, ok := f.getStatus(conn, status.Notification)
			if !ok {
				continue
			}

			switch status.Notification.TypeId {
			case agency.Notification_ACTION_NEEDED:
				f.handleAction(job, status.Notification, protocolStatus)
			case agency.Notification_STATUS_UPDATE:
				f.handleStatus(conn, job, status.Notification, protocolStatus)
			case agency.Notification_ANSWER_NEEDED_PING:
			case agency.Notification_ANSWER_NEEDED_ISSUE_PROPOSE:
			case agency.Notification_ANSWER_NEEDED_PROOF_PROPOSE:
			case agency.Notification_ANSWER_NEEDED_PROOF_VERIFY:
				// TODO?
			}
		}
	}()

	return err
}

func (f *Agency) releaseCompleted(
	conn client.Conn,
	job *model.JobInfo,
	notification *agency.Notification,
	status *agency.ProtocolStatus,
) {
	defer err2.Catch(func(err error) {
		glog.Errorf("Failure when releasing protocol: %s", err.Error())
	})
	cmd := agency.NewDIDCommClient(conn)
	_, err := cmd.Release(f.ctx, status.State.ProtocolId)
	err2.Check(err)
}
