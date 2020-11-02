package resolver

import (
	"context"

	"github.com/findy-network/findy-agent-vault/agency"

	"github.com/findy-network/findy-agent-vault/graph/model"
	"github.com/golang/glog"
	"github.com/lainio/err2"
)

func (r *mutationResolver) Connect(_ context.Context, input model.ConnectInput) (res *model.Response, err error) {
	defer err2.Return(&err)
	glog.V(logLevelMedium).Info("mutationResolver:Connect")

	_, err = agency.Instance.Connect(input.Invitation)
	err2.Check(err)

	res = &model.Response{Ok: true}

	addEvent("Sent connection request", model.ProtocolTypeConnection, "")

	return
}
