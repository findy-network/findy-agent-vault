package invitation

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"image/png"

	agency "github.com/findy-network/findy-agent-vault/agency/model"
	"github.com/findy-network/findy-agent-vault/graph/model"
	didexchange "github.com/findy-network/findy-common-go/std/didexchange/invitation"
	"github.com/lainio/err2"
	"github.com/lainio/err2/try"
	"github.com/skip2/go-qrcode"
)

const imageSize = 256

func FromURLParam(raw string) (res *model.InvitationResponse, err error) {
	err2.Return(&err)

	qrCode := try.To1(strToQRCode(raw))

	invitation := didexchange.Invitation{}
	// ignore error on purpose
	_ = json.Unmarshal([]byte(raw), &invitation)

	// TODO: parse invitation in URL-format

	res = &model.InvitationResponse{
		ID:       invitation.ID,
		Endpoint: invitation.ServiceEndpoint,
		Label:    invitation.Label,
		Raw:      raw,
		ImageB64: qrCode,
	}

	return
}

func FromAgency(data *agency.InvitationData) (res *model.InvitationResponse, err error) {
	err2.Return(&err)

	qrCode := try.To1(strToQRCode(data.Raw))

	res = &model.InvitationResponse{
		ID:       data.Data.ID,
		Endpoint: data.Data.ServiceEndpoint,
		Label:    data.Data.Label,
		Raw:      data.Raw,
		ImageB64: qrCode,
	}

	return
}

func strToQRCode(str string) (res string, err error) {
	defer err2.Return(&err)

	code := try.To1(qrcode.New(str, qrcode.Low))

	var buf bytes.Buffer
	try.To(png.Encode(&buf, code.Image(imageSize)))

	res = base64.StdEncoding.EncodeToString(buf.Bytes())
	return
}
