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

func (f *Agency) getStatus(conn client.Conn, notification *agency.Notification) (*agency.ProtocolStatus, error) {
	ctx := context.Background()
	didComm := agency.NewDIDCommClient(conn)
	return didComm.Status(ctx, &agency.ProtocolID{
		TypeId:           notification.ProtocolType,
		Id:               notification.ProtocolId,
		NotificationTime: notification.Timestamp,
	})
}

func (f *Agency) handleStatus(
	job *model.JobInfo, notification *agency.Notification,
	status *agency.ProtocolStatus,
) {
	now := utils.CurrentTimeMs()
	switch notification.ProtocolType {
	case agency.Protocol_CONNECT:
		connection := status.GetConnection()
		f.vault.AddConnection(
			job,
			connection.MyDid,
			connection.TheirDid,
			connection.TheirEndpoint,
			connection.TheirLabel,
		)
	case agency.Protocol_BASIC_MESSAGE:
		message := status.GetBasicMessage()
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
		// TODO
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

	conn := userListenClient(a)

	ch, err := conn.Listen(f.ctx, &agency.ClientID{Id: a.TenantID})
	err2.Check(err)

	go func() {
		for {
			status, ok := <-ch
			if !ok {
				glog.Warningln("closed from server")
				conn.Close()
				break
			}
			utils.LogLow().Infoln("received notification:",
				status.Notification.TypeId,
				status.Notification.Role,
				status.Notification.ProtocolId)

			job := &model.JobInfo{
				TenantID:     a.TenantID,
				JobID:        status.Notification.ProtocolId,
				ConnectionID: status.Notification.ConnectionId,
			}

			protocolStatus, statusErr := f.getStatus(conn, status.Notification)
			if err != nil {
				glog.Error(statusErr)
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
