// +build !findy
// +build !findy_grpc

package agency

import (
	"encoding/json"
	"time"

	"github.com/bxcodec/faker"
	generator "github.com/findy-network/findy-agent-vault/tools/faker"
	"github.com/google/uuid"
	"github.com/lainio/err2"
)

type Mock struct {
	listener Listener
}

type invitation struct {
	ServiceEndpoint string   `json:"serviceEndpoint,omitempty" faker:"url"`
	RecipientKeys   []string `json:"recipientKeys,omitempty" faker:"-"`
	ID              string   `json:"@id,omitempty" faker:"uuid_hyphenated"`
	Label           string   `json:"label,omitempty" faker:"first_name"`
	Type            string   `json:"@type,omitempty" faker:"-"`
}

var Instance Agency = &Mock{}

func (m *Mock) Init(l Listener) {
	m.listener = l
}

func (m *Mock) Invite(a *Agent) (result, id string, err error) {
	defer err2.Return(&err)

	inv := invitation{}
	err = faker.FakeData(&inv)
	err2.Check(err)

	inv.RecipientKeys = append(inv.RecipientKeys, "CDdVp7CyP9Ued38FpFd8rqxF3eEKhrnjAsPWf6LEeLJC")
	inv.Type = "did:sov:BzCbsNYhMrjHiqZDTUASHg;spec/connections/1.0/invitation"

	jsonBytes := err2.Bytes.Try(json.Marshal(&inv))
	result = string(jsonBytes)
	id = inv.ID

	return
}

func (m *Mock) Connect(a *Agent, strInvitation string) (id string, err error) {
	defer err2.Return(&err)

	inv := invitation{}
	err2.Check(json.Unmarshal([]byte(strInvitation), &inv))

	id = inv.ID

	job := &JobInfo{TenantID: a.TenantID, JobID: id, ConnectionID: id}

	time.AfterFunc(time.Second, func() {
		if connections, err := generator.FakeConnections(1, true); err == nil {
			connection := connections[0]
			m.listener.AddConnection(job, connection.OurDid, connection.TheirDid, connection.TheirEndpoint, connection.TheirLabel)
		}
	})

	return
}

func (m *Mock) SendMessage(a *Agent, connectionID, message string) (id string, err error) {
	defer err2.Return(&err)

	id = uuid.New().String()

	job := &JobInfo{TenantID: a.TenantID, JobID: id, ConnectionID: connectionID}

	m.listener.AddMessage(job, message, true)
	time.AfterFunc(time.Second, func() {
		if messages, err := generator.FakeMessages(1); err == nil {
			msg := messages[0]
			// reply
			m.listener.AddMessage(job, msg.Message, false)
		}
	})

	return
}

func (m *Mock) ResumeCredentialOffer(a *Agent, id string, accept bool) (err error) {
	return
}

func (m *Mock) ResumeProofRequest(a *Agent, id string, accept bool) (err error) {
	return
}
