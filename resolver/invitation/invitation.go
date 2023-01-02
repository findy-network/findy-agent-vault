package invitation

import (
	"bytes"
	"encoding/base64"
	"image/png"

	agency "github.com/findy-network/findy-agent-vault/agency/model"
	"github.com/findy-network/findy-agent-vault/graph/model"
	"github.com/findy-network/findy-common-go/std/didexchange/invitation"
	"github.com/lainio/err2"
	"github.com/lainio/err2/try"
	"github.com/skip2/go-qrcode"
)

const imageSize = 256

func FromURLParam(raw string) (res *model.InvitationResponse, err error) {
	defer err2.Handle(&err)

	qrCode := try.To1(strToQRCode(raw))

	inv := try.To1(invitation.Translate(raw))

	res = &model.InvitationResponse{
		ID:       inv.ID(),
		Endpoint: inv.Services()[0].ServiceEndpoint,
		Label:    inv.Label(),
		Raw:      raw,
		ImageB64: qrCode,
	}

	return
}

func FromAgency(data *agency.InvitationData) (res *model.InvitationResponse, err error) {
	defer err2.Handle(&err)

	qrCode := try.To1(strToQRCode(data.Raw))

	inv := try.To1(invitation.Translate(data.Raw))

	res = &model.InvitationResponse{
		ID:       inv.ID(),
		Endpoint: inv.Services()[0].ServiceEndpoint,
		Label:    inv.Label(),
		Raw:      data.Raw,
		ImageB64: qrCode,
	}

	return
}

func strToQRCode(str string) (res string, err error) {
	defer err2.Handle(&err)

	code := try.To1(qrcode.New(str, qrcode.Low))

	var buf bytes.Buffer
	try.To(png.Encode(&buf, code.Image(imageSize)))

	res = base64.StdEncoding.EncodeToString(buf.Bytes())
	return
}
