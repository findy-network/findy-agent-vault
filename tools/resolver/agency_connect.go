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

func (r *mutationResolver) Connect(_ context.Context, input model.ConnectInput) (res *model.Response, err error) {
	defer err2.Return(&err)
	glog.V(logLevelMedium).Info("mutationResolver:Connect")

	id, err := agency.Instance.Connect(input.Invitation)
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
		ID:            id,
		OurDid:        ourDID,
		TheirDid:      theirDID,
		TheirEndpoint: theirEndpoint,
		TheirLabel:    theirLabel,
		ApprovedMs:    currentTime,
		CreatedMs:     currentTime,
	})
}
