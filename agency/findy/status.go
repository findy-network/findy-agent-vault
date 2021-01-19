package findy

import (
	"context"

	"github.com/findy-network/findy-agent-api/grpc/agency"
	"github.com/findy-network/findy-agent-vault/agency/model"
	graph "github.com/findy-network/findy-agent-vault/graph/model"
	"github.com/findy-network/findy-agent-vault/utils"
	"github.com/findy-network/findy-grpc/agency/client"
	"github.com/golang/glog"
	"github.com/lainio/err2"
)

func (f *Agency) userListenClient(a *model.Agent) client.Conn {
	config := client.BuildClientConnBase(f.tlsPath, agencyHost, agencyPort, f.options)
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

func (f *Agency) handleStatus(
	job *model.JobInfo, notification *agency.Notification,
	status *agency.ProtocolStatus,
) {
	// TODO: check status (failed/successful?)

	now := utils.CurrentTimeMs()
	switch notification.ProtocolType {
	case agency.Protocol_CONNECT:
		connection := status.GetConnection()
		if connection == nil {
			glog.Errorf("Received invalid connection object for %s", job.JobID)
			return
		}

		f.vault.AddConnection(
			job,
			connection.MyDid,
			connection.TheirDid,
			connection.TheirEndpoint,
			connection.TheirLabel,
		)

	case agency.Protocol_BASIC_MESSAGE:
		message := status.GetBasicMessage()
		if message == nil {
			glog.Errorf("Received invalid message object for %s", job.JobID)
			return
		}

		f.vault.AddMessage(
			job,
			message.Content,
			message.SentByMe,
			// TODO: delivered?
		)
	case agency.Protocol_ISSUE:
		f.vault.UpdateCredential(
			job,
			nil,
			&now,
			nil,
		)
	case agency.Protocol_PROOF:
		f.vault.UpdateProof(
			job,
			nil,
			&now,
			nil,
		)
	case agency.Protocol_NONE:
	case agency.Protocol_TRUST_PING:
	}
}

func (f *Agency) handleAction(
	job *model.JobInfo,
	notification *agency.Notification,
	status *agency.ProtocolStatus,
) {
	switch notification.ProtocolType {
	case agency.Protocol_ISSUE:
		credential := status.GetIssue()
		if credential == nil {
			glog.Errorf("Received invalid credential issue object for %s", job.JobID)
			return
		}

		role := graph.CredentialRoleHolder
		if notification.Role != agency.Protocol_ADDRESSEE {
			role = graph.CredentialRoleIssuer
		}
		values := make([]*graph.CredentialValue, 0)
		for _, v := range credential.Attrs {
			values = append(values, &graph.CredentialValue{
				Name:  v.Name,
				Value: v.Value,
			})
		}
		// TODO: what if we are issuer?
		f.vault.AddCredential(job, role, credential.SchemaId, credential.CredDefId, values, false)
	case agency.Protocol_PROOF:
		proof := status.GetProof()
		if proof == nil {
			glog.Errorf("Received invalid proof object for %s", job.JobID)
			return
		}

		role := graph.ProofRoleProver
		if notification.Role != agency.Protocol_ADDRESSEE {
			role = graph.ProofRoleVerifier
		}
		attributes := make([]*graph.ProofAttribute, 0)
		for _, v := range proof.Attrs {
			value := "" // TODO: get also values from notification?
			attributes = append(attributes, &graph.ProofAttribute{
				Name:      v.Name,
				Value:     &value,
				CredDefID: v.CredDefId,
			})
		}
		// TODO: what if we are verifier?
		f.vault.AddProof(job, role, attributes, false)
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

		// TODO: fail job if error happens?
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
				f.handleStatus(job, status.Notification, protocolStatus)
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
