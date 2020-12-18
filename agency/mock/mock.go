package mock

import (
	"encoding/json"
	"time"

	"github.com/findy-network/findy-agent-vault/agency/model"
	"github.com/findy-network/findy-agent-vault/db/fake"

	"github.com/bxcodec/faker"
	"github.com/google/uuid"
	"github.com/lainio/err2"
)

type Mock struct {
	listener model.Listener
}

func (m *Mock) Init(l model.Listener, agents []*model.Agent) {
	m.listener = l
}

func (m *Mock) Invite(a *model.Agent) (result, id string, err error) {
	defer err2.Return(&err)

	inv := model.Invitation{}
	err = faker.FakeData(&inv)
	err2.Check(err)

	inv.RecipientKeys = append(inv.RecipientKeys, "CDdVp7CyP9Ued38FpFd8rqxF3eEKhrnjAsPWf6LEeLJC")
	inv.Type = "did:sov:BzCbsNYhMrjHiqZDTUASHg;spec/connections/1.0/invitation"

	jsonBytes := err2.Bytes.Try(json.Marshal(&inv))
	result = string(jsonBytes)
	id = inv.ID

	return
}

func (m *Mock) Connect(a *model.Agent, strInvitation string) (id string, err error) {
	defer err2.Return(&err)

	inv := model.Invitation{}
	err2.Check(json.Unmarshal([]byte(strInvitation), &inv))

	id = inv.ID

	job := &model.JobInfo{TenantID: a.TenantID, JobID: id, ConnectionID: id}

	time.AfterFunc(time.Second, func() {
		connection := fake.Connection(a.TenantID)
		m.listener.AddConnection(job, connection.OurDid, connection.TheirDid, connection.TheirEndpoint, connection.TheirLabel)
	})

	return
}

func (m *Mock) SendMessage(a *model.Agent, connectionID, message string) (id string, err error) {
	defer err2.Return(&err)

	id = uuid.New().String()

	job := &model.JobInfo{TenantID: a.TenantID, JobID: id, ConnectionID: connectionID}

	m.listener.AddMessage(job, message, true)
	time.AfterFunc(time.Second, func() {
		message := fake.Message(a.TenantID, connectionID)
		m.listener.AddMessage(job, message.Message, false)
	})
	return
}

func (m *Mock) ResumeCredentialOffer(a *model.Agent, id string, accept bool) (err error) {
	return
}

func (m *Mock) ResumeProofRequest(a *model.Agent, id string, accept bool) (err error) {
	return
}
