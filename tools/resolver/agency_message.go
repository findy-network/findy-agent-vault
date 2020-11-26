package resolver

import (
	"context"

	data "github.com/findy-network/findy-agent-vault/tools/data/model"
	"github.com/findy-network/findy-agent-vault/tools/utils"

	"github.com/findy-network/findy-agent-vault/agency"
	"github.com/findy-network/findy-agent-vault/graph/model"
	"github.com/golang/glog"
	"github.com/lainio/err2"
)

func (r *mutationResolver) SendMessage(ctx context.Context, input model.MessageInput) (res *model.Response, err error) {
	defer err2.Return(&err)
	glog.V(logLevelMedium).Info("mutationResolver:SendMessage")

	id, err := agency.Instance.SendMessage(ctx, input.ConnectionID, input.Message)
	err2.Check(err)

	res = &model.Response{Ok: true}

	addJob(
		id,
		model.ProtocolTypeBasicMessage,
		&id,
		true,
		&input.ConnectionID,
		"Sent basic message")

	return
}

func (l *agencyListener) AddMessage(connectionID, id, message string, sentByMe bool) {
	currentTime := utils.CurrentTimeMs()
	msg := data.InternalMessage{
		BaseObject: &data.BaseObject{
			ID:        id,
			CreatedMs: currentTime,
		},
		Message:    message,
		PairwiseID: connectionID,
		SentByMe:   sentByMe,
		Delivered:  nil,
	}
	desc := msg.Description()
	state.Messages.Append(&msg)
	glog.Infof("Added message %s for connection %s", id, connectionID)

	addJobWithStatus(
		id,
		model.ProtocolTypeBasicMessage,
		&id,
		sentByMe,
		&connectionID,
		desc,
		model.JobStatusComplete,
		model.JobResultSuccess)
}

func (l *agencyListener) UpdateMessage(connectionID, id, delivered bool) {
	// TODO complete message job
}
