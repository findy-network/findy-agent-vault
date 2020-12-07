package utils

import (
	"bytes"
	"encoding/base64"
	"image/png"

	"github.com/lainio/err2"

	qrcode "github.com/skip2/go-qrcode"
)

const imageSize = 256

func StrToQRCode(str string) (res string, err error) {
	defer err2.Return(&err)

	code, err := qrcode.New(str, qrcode.Low)
	err2.Check(err)

	var buf bytes.Buffer
	err = png.Encode(&buf, code.Image(imageSize))
	err2.Check(err)

	res = base64.StdEncoding.EncodeToString(buf.Bytes())
	return
}