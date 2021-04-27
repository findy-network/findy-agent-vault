package utils

import (
	"encoding/json"

	"github.com/findy-network/findy-agent-vault/graph/model"
	"github.com/findy-network/findy-common-go/std/didexchange/invitation"
	"github.com/lainio/err2"
)

func FromAriesInvitation(invitationStr string) (res *model.InvitationResponse, err error) {
	err2.Return(&err)

	inv := invitation.Invitation{}
	err2.Check(json.Unmarshal([]byte(invitationStr), &inv))

	qrCode, err := StrToQRCode(invitationStr)
	err2.Check(err)

	res = &model.InvitationResponse{
		ID:       inv.ID,
		Endpoint: inv.ServiceEndpoint,
		Label:    inv.Label,
		Raw:      invitationStr,
		ImageB64: qrCode,
	}

	return
}
