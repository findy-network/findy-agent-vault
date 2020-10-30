// +build !findy

package agency

import (
	"encoding/json"

	"github.com/bxcodec/faker"
	"github.com/lainio/err2"
)

type Mock struct{}

type invitation struct {
	ServiceEndpoint string   `json:"serviceEndpoint,omitempty" faker:"url"`
	RecipientKeys   []string `json:"recipientKeys,omitempty" faker:"-"`
	ID              string   `json:"@id,omitempty" faker:"uuid_hyphenated"`
	Label           string   `json:"label,omitempty" faker:"first_name"`
	Type            string   `json:"@type,omitempty" faker:"-"` //did:sov:BzCbsNYhMrjHiqZDTUASHg;spec/connections/1.0/invitation
}

var Instance Agency = &Mock{}

func (m *Mock) Init() {}

func (m *Mock) Invite() (result string, err error) {
	defer err2.Return(&err)

	inv := invitation{}
	err = faker.FakeData(&inv)
	err2.Check(err)

	inv.RecipientKeys = append(inv.RecipientKeys, "CDdVp7CyP9Ued38FpFd8rqxF3eEKhrnjAsPWf6LEeLJC")
	inv.Type = "did:sov:BzCbsNYhMrjHiqZDTUASHg;spec/connections/1.0/invitation"

	jsonBytes := err2.Bytes.Try(json.Marshal(&inv))
	result = string(jsonBytes)

	return
}

func (m *Mock) Connect() (string, error) {
	return "", nil
}
