package fake

import (
	crand "crypto/rand"
	"math/big"
	"reflect"
	"time"

	"github.com/bxcodec/faker/v3"
	"github.com/findy-network/findy-agent-vault/db/model"
	"github.com/findy-network/findy-agent-vault/db/store"
	graph "github.com/findy-network/findy-agent-vault/graph/model"
	"github.com/findy-network/findy-agent-vault/utils"
	"github.com/lainio/err2"
)

const FakeCloudDID = "cloudDID"

func random(n int) int {
	val, err := crand.Int(crand.Reader, big.NewInt(int64(n)))
	if err != nil {
		panic(err)
	}
	return int(val.Int64())
}

func AddMessages(db store.DB, tenantID, connectionID string, count int) []*model.Message {
	messages := make([]*model.Message, count)
	for i := 0; i < count; i++ {
		message := Message(tenantID, connectionID)
		messages[i] = message
	}

	newMessages := make([]*model.Message, count)
	for index, message := range messages {
		c, err := db.AddMessage(message)
		err2.Check(err)
		time.Sleep(time.Millisecond) // generate different timestamps for items

		newMessages[index] = c
	}

	utils.LogTrace().Infof("Generated %d messages for tenant %s", len(newMessages), tenantID)

	return newMessages
}

func AddEvents(db store.DB, tenantID, connectionID string, jobID *string, count int) []*model.Event {
	events := make([]*model.Event, count)
	for i := 0; i < count; i++ {
		event := fakeEvent(tenantID, connectionID, jobID)
		events[i] = event
	}

	newEvents := make([]*model.Event, count)
	for index, event := range events {
		c, err := db.AddEvent(event)
		err2.Check(err)
		time.Sleep(time.Millisecond) // generate different timestamps for items

		newEvents[index] = c
	}

	utils.LogTrace().Infof("Generated %d events for tenant %s", len(newEvents), tenantID)

	return newEvents
}

func AddJobs(db store.DB, tenantID, connectionID string, count int) []*model.Job {
	return addJobs(db, tenantID, connectionID, nil, nil, nil, nil, count, graph.JobStatusComplete)
}

func AddConnectionJobs(db store.DB, tenantID, connectionID, protocolConnectionID string, count int) []*model.Job {
	return addJobs(db, tenantID, connectionID, &protocolConnectionID, nil, nil, nil, count, graph.JobStatusComplete)
}

func AddCredentialJobs(db store.DB, tenantID, connectionID, protocolCredentialID string, count int) []*model.Job {
	return addJobs(db, tenantID, connectionID, nil, &protocolCredentialID, nil, nil, count, graph.JobStatusComplete)
}

func AddProofJobs(db store.DB, tenantID, connectionID, protocolProofID string, count int, status graph.JobStatus) []*model.Job {
	return addJobs(db, tenantID, connectionID, nil, nil, &protocolProofID, nil, count, status)
}

func AddMessageJobs(db store.DB, tenantID, connectionID, protocolMessageID string, count int) []*model.Job {
	return addJobs(db, tenantID, connectionID, nil, nil, nil, &protocolMessageID, count, graph.JobStatusComplete)
}

func AddCredentials(db store.DB, tenantID, connectionID string, count int) []*model.Credential {
	credentials := make([]*model.Credential, count)
	for i := 0; i < count; i++ {
		credential := Credential(tenantID, connectionID)
		credentials[i] = credential
	}

	newCredentials := make([]*model.Credential, count)
	for index, credential := range credentials {
		c, err := db.AddCredential(credential)
		err2.Check(err)
		time.Sleep(time.Millisecond) // generate different timestamps for items

		now := time.Now().UTC()
		c.Approved = now
		c.Issued = now
		_, err = db.UpdateCredential(c)
		err2.Check(err)

		newCredentials[index] = c
	}

	utils.LogTrace().Infof("Generated %d credentials for tenant %s", len(newCredentials), tenantID)

	return newCredentials
}

func AddProofs(db store.DB, tenantID, connectionID string, count int, verify bool) []*model.Proof {
	proofs := make([]*model.Proof, count)
	for i := 0; i < count; i++ {
		proof := Proof(tenantID, connectionID)
		proofs[i] = proof
	}

	newProofs := make([]*model.Proof, count)
	for index, proof := range proofs {
		p, err := db.AddProof(proof)
		err2.Check(err)
		time.Sleep(time.Millisecond) // generate different timestamps for items

		now := time.Now().UTC()
		if verify {
			p.Approved = &now
			p.Verified = &now
			_, err = db.UpdateProof(p)
			err2.Check(err)
		}

		newProofs[index] = p
	}

	utils.LogTrace().Infof("Generated %d proofs for tenant %s", len(newProofs), tenantID)

	return newProofs
}

func AddConnections(db store.DB, tenantID string, count int) []*model.Connection {
	connections := make([]*model.Connection, count)
	for i := 0; i < count; i++ {
		connection := Connection(tenantID)
		connections[i] = connection
	}

	newConnections := make([]*model.Connection, count)
	for index, connection := range connections {
		c, err := db.AddConnection(connection)
		err2.Check(err)
		time.Sleep(time.Millisecond) // generate different timestamps for items
		newConnections[index] = c
	}

	utils.LogTrace().Infof("Generated %d connections for tenant %s", len(newConnections), tenantID)

	return newConnections
}

func AddAgent(db store.DB) *model.Agent {
	_ = faker.AddProvider("agentId", func(v reflect.Value) (interface{}, error) {
		return FakeCloudDID, nil
	})
	var err error
	agent := fakeAgent()

	agent, err = db.AddAgent(agent)
	err2.Check(err)

	utils.LogTrace().Infof("Generated tenant %s with agent id %s", agent.ID, agent.AgentID)
	return agent
}

func AddData(db store.DB) {
	agent := AddAgent(db)

	count := 5
	eventFactor := 3

	connections := AddConnections(db, agent.ID, count)
	AddCredentials(db, agent.ID, connections[0].ID, count)
	AddProofs(db, agent.ID, connections[0].ID, count, true)
	AddMessages(db, agent.ID, connections[0].ID, count)
	jobs := AddJobs(db, agent.ID, connections[0].ID, count)
	AddEvents(db, agent.ID, connections[0].ID, &jobs[0].ID, count*eventFactor)
}

func addJobs(
	db store.DB,
	tenantID, connectionID string,
	protocolConnectionID, protocolCredentialID, protocolProofID, protocolMessageID *string,
	count int,
	status graph.JobStatus,
) []*model.Job {
	jobs := make([]*model.Job, count)
	for i := 0; i < count; i++ {
		job := fakeJob(
			faker.UUIDHyphenated(),
			tenantID,
			connectionID,
			protocolConnectionID,
			protocolCredentialID,
			protocolProofID,
			protocolMessageID,
			status,
		)
		jobs[i] = job
	}

	newJobs := make([]*model.Job, count)
	for index, job := range jobs {
		c, err := db.AddJob(job)
		err2.Check(err)
		time.Sleep(time.Millisecond) // generate different timestamps for items

		newJobs[index] = c
	}

	utils.LogTrace().Infof("Generated %d jobs for tenant %s", len(newJobs), tenantID)

	return newJobs
}

func fakeAgent() *model.Agent {
	agent := &model.Agent{}
	err2.Check(faker.FakeData(&agent))
	return agent
}

func Connection(tenantID string) *model.Connection {
	_ = faker.AddProvider("organisationLabel", func(v reflect.Value) (interface{}, error) {
		orgs := []string{"Bank", "Ltd", "Agency", "Company", "United"}
		index := random(len(orgs))
		return faker.LastName() + " " + orgs[index], nil
	})

	connection := &model.Connection{}
	err2.Check(faker.FakeData(connection))
	connection.ID = faker.UUIDHyphenated()
	connection.TenantID = tenantID
	return connection
}

func Credential(tenantID, connectionID string) *model.Credential {
	_ = faker.AddProvider("credentialAttributes", func(v reflect.Value) (interface{}, error) {
		return []*graph.CredentialValue{
			{Name: "name1", Value: "value1"},
			{Name: "name2", Value: "value2"},
			{Name: "name3", Value: "value3"},
		}, nil
	})
	credential := &model.Credential{}
	err2.Check(faker.FakeData(credential))
	credential.TenantID = tenantID
	credential.ConnectionID = connectionID
	return credential
}

func Proof(tenantID, connectionID string) *model.Proof {
	_ = faker.AddProvider("proofAttributes", func(v reflect.Value) (interface{}, error) {
		return []*graph.ProofAttribute{
			{Name: "name1", CredDefID: "credDefId1"},
			{Name: "name2", CredDefID: "credDefId2"},
			{Name: "name3", CredDefID: "credDefId3"},
		}, nil
	})
	proof := model.NewProof("", nil)
	err2.Check(faker.FakeData(proof))
	proof = model.NewProof(tenantID, proof)
	proof.TenantID = tenantID
	proof.ConnectionID = connectionID
	return proof
}

func fakeEvent(tenantID, connectionID string, jobID *string) *model.Event {
	event := &model.Event{}
	err2.Check(faker.FakeData(event))
	event.TenantID = tenantID
	event.ConnectionID = &connectionID
	event.JobID = jobID
	return event
}

func fakeJob(
	id, tenantID, connectionID string,
	protocolConnectionID, protocolCredentialID, protocolProofID, protocolMessageID *string,
	status graph.JobStatus,
) *model.Job {
	job := &model.Job{}
	err2.Check(faker.FakeData(job))
	job.ID = id
	job.TenantID = tenantID
	job.ConnectionID = &connectionID
	job.ProtocolConnectionID = protocolConnectionID
	job.ProtocolCredentialID = protocolCredentialID
	if job.ProtocolCredentialID != nil {
		job.ProtocolType = graph.ProtocolTypeCredential
	}
	job.ProtocolProofID = protocolProofID
	if job.ProtocolProofID != nil {
		job.ProtocolType = graph.ProtocolTypeProof
	}
	job.ProtocolMessageID = protocolMessageID
	if job.ProtocolMessageID != nil {
		job.ProtocolType = graph.ProtocolTypeBasicMessage
	}
	job.Status = status
	return job
}

func Message(tenantID, connectionID string) *model.Message {
	message := model.NewMessage("", nil)
	err2.Check(faker.FakeData(message))
	message = model.NewMessage(tenantID, message)
	message.TenantID = tenantID
	message.ConnectionID = connectionID
	return message
}
