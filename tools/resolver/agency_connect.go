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

func (r *mutationResolver) Connect(ctx context.Context, input model.ConnectInput) (res *model.Response, err error) {
	defer err2.Return(&err)
	glog.V(logLevelMedium).Info("mutationResolver:Connect")

	id, err := agency.Instance.Connect(ctx, input.Invitation)
	err2.Check(err)

	res = &model.Response{Ok: true}

	addJob(
		id,
		model.ProtocolTypeConnection,
		nil,
		false,
		nil,
		"Sent connection request")

	return
}

func (l *agencyListener) AddConnection(id, ourDID, theirDID, theirEndpoint, theirLabel string) {
	currentTime := utils.CurrentTimeMs()
	doAddConnection(&data.InternalPairwise{
		BaseObject: &data.BaseObject{
			ID:        id,
			CreatedMs: currentTime,
		},
		OurDid:        ourDID,
		TheirDid:      theirDID,
		TheirEndpoint: theirEndpoint,
		TheirLabel:    theirLabel,
		ApprovedMs:    currentTime,
	})
}

func doAddConnection(connection *data.InternalPairwise) {
	items := state.Connections().Objects()
	connection.CreatedMs = utils.CurrentTimeMs()
	initiatedByUs := state.Jobs.IsJobInitiatedByUs(connection.ID)
	if initiatedByUs != nil {
		connection.Invited = *initiatedByUs
	}
	items.Append(connection)
	glog.Infof("Added connection %s", connection.ID)
	updateJob(
		connection.ID,
		&connection.ID,
		&connection.ID,
		model.JobStatusComplete,
		model.JobResultSuccess,
		"Established connection to "+connection.TheirLabel)
}
