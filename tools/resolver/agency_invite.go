package resolver

import (
	"context"

	"github.com/findy-network/findy-agent-vault/utils"

	"github.com/findy-network/findy-agent-vault/agency"
	"github.com/golang/glog"
	"github.com/lainio/err2"

	"github.com/findy-network/findy-agent-vault/graph/model"
)

func (r *mutationResolver) Invite(ctx context.Context) (resp *model.InvitationResponse, err error) {
	defer err2.Return(&err)
	glog.V(logLevelMedium).Info("mutationResolver:Invite")

	str, err := agency.Instance.Invite()
	err2.Check(err)

	img, err := utils.StrToQRCode(str)
	err2.Check(err)

	resp = &model.InvitationResponse{
		Invitation: str,
		ImageB64:   img,
	}

	addEvent("Created new invitation", model.ProtocolTypeConnection)

	return
}
