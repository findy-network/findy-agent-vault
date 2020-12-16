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
		message := fakeMessage(tenantID, connectionID)
		messages[i] = message
	}

	newMessages := make([]*model.Message, count)
	for index, message := range messages {
		c, err := db.AddMessage(message)
		err2.Check(err)
		time.Sleep(time.Millisecond) // generate different timestamps for items

		newMessages[index] = c
	}

	utils.LogMed().Infof("Generated %d messages for tenant %s", len(newMessages), tenantID)

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

	utils.LogMed().Infof("Generated %d events for tenant %s", len(newEvents), tenantID)

	return newEvents
}

func AddJobs(db store.DB, tenantID, connectionID string, count int) []*model.Job {
	return addJobs(db, tenantID, connectionID, nil, nil, nil, nil, count)
}

func AddConnectionJobs(db store.DB, tenantID, connectionID, protocolConnectionID string, count int) []*model.Job {
	return addJobs(db, tenantID, connectionID, &protocolConnectionID, nil, nil, nil, count)
}

func AddCredentialJobs(db store.DB, tenantID, connectionID, protocolCredentialID string, count int) []*model.Job {
	return addJobs(db, tenantID, connectionID, nil, &protocolCredentialID, nil, nil, count)
}

func AddProofJobs(db store.DB, tenantID, connectionID, protocolProofID string, count int) []*model.Job {
	return addJobs(db, tenantID, connectionID, nil, nil, &protocolProofID, nil, count)
}

func AddMessageJobs(db store.DB, tenantID, connectionID, protocolMessageID string, count int) []*model.Job {
	return addJobs(db, tenantID, connectionID, nil, nil, nil, &protocolMessageID, count)
}

func AddCredentials(db store.DB, tenantID, connectionID string, count int) []*model.Credential {
	_ = faker.AddProvider("credentialAttributes", func(v reflect.Value) (interface{}, error) {
		return []*graph.CredentialValue{
			{Name: "name1", Value: "value1"},
			{Name: "name2", Value: "value2"},
			{Name: "name3", Value: "value3"},
		}, nil
	})

	credentials := make([]*model.Credential, count)
	for i := 0; i < count; i++ {
		credential := fakeCredential(tenantID, connectionID)
		credentials[i] = credential
	}

	newCredentials := make([]*model.Credential, count)
	for index, credential := range credentials {
		c, err := db.AddCredential(credential)
		err2.Check(err)
		time.Sleep(time.Millisecond) // generate different timestamps for items

		now := time.Now().UTC()
		c.Approved = &now
		c.Issued = &now
		_, err = db.UpdateCredential(c)
		err2.Check(err)

		newCredentials[index] = c
	}

	utils.LogMed().Infof("Generated %d credentials for tenant %s", len(newCredentials), tenantID)

	return newCredentials
}

func AddProofs(db store.DB, tenantID, connectionID string, count int) []*model.Proof {
	_ = faker.AddProvider("proofAttributes", func(v reflect.Value) (interface{}, error) {
		value := "value"
		return []*graph.ProofAttribute{
			{Name: "name1", Value: &value, CredDefID: "credDefId1"},
			{Name: "name2", Value: &value, CredDefID: "credDefId2"},
			{Name: "name3", Value: &value, CredDefID: "credDefId3"},
		}, nil
	})

	proofs := make([]*model.Proof, count)
	for i := 0; i < count; i++ {
		proof := fakeProof(tenantID, connectionID)
		proofs[i] = proof
	}

	newProofs := make([]*model.Proof, count)
	for index, proof := range proofs {
		p, err := db.AddProof(proof)
		err2.Check(err)
		time.Sleep(time.Millisecond) // generate different timestamps for items

		now := time.Now().UTC()
		p.Approved = &now
		p.Verified = &now
		_, err = db.UpdateProof(p)
		err2.Check(err)

		newProofs[index] = p
	}

	utils.LogMed().Infof("Generated %d proofs for tenant %s", len(newProofs), tenantID)

	return newProofs
}

func AddConnections(db store.DB, tenantID string, count int) []*model.Connection {
	_ = faker.AddProvider("organisationLabel", func(v reflect.Value) (interface{}, error) {
		orgs := []string{"Bank", "Ltd", "Agency", "Company", "United"}
		index := random(len(orgs))
		return faker.LastName() + " " + orgs[index], nil
	})

	connections := make([]*model.Connection, count)
	for i := 0; i < count; i++ {
		connection := fakeConnection(tenantID)
		connections[i] = connection
	}

	newConnections := make([]*model.Connection, count)
	for index, connection := range connections {
		c, err := db.AddConnection(connection)
		err2.Check(err)
		time.Sleep(time.Millisecond) // generate different timestamps for items
		newConnections[index] = c
	}

	utils.LogMed().Infof("Generated %d connections for tenant %s", len(newConnections), tenantID)

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

	utils.LogMed().Infof("Generated tenant %s with agent id %s", agent.ID, agent.AgentID)
	return agent
}

func AddData(db store.DB) {
	agent := AddAgent(db)

	AddConnections(db, agent.ID, 5)
}

func addJobs(
	db store.DB,
	tenantID, connectionID string,
	protocolConnectionID, protocolCredentialID, protocolProofID, protocolMessageID *string,
	count int,
) []*model.Job {
	jobs := make([]*model.Job, count)
	for i := 0; i < count; i++ {
		job := fakeJob(
			tenantID,
			connectionID,
			protocolConnectionID,
			protocolCredentialID,
			protocolProofID,
			protocolMessageID,
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

	utils.LogMed().Infof("Generated %d jobs for tenant %s", len(newJobs), tenantID)

	return newJobs
}

func fakeAgent() *model.Agent {
	agent := model.NewAgent(nil)
	err2.Check(faker.FakeData(&agent))
	return model.NewAgent(agent)
}

func fakeConnection(tenantID string) *model.Connection {
	connection := model.NewConnection(nil)
	err2.Check(faker.FakeData(connection))
	connection = model.NewConnection(connection)
	connection.TenantID = tenantID
	return connection
}

func fakeCredential(tenantID, connectionID string) *model.Credential {
	credential := model.NewCredential(nil)
	err2.Check(faker.FakeData(credential))
	credential = model.NewCredential(credential)
	credential.TenantID = tenantID
	credential.ConnectionID = connectionID
	return credential
}

func fakeProof(tenantID, connectionID string) *model.Proof {
	proof := model.NewProof(nil)
	err2.Check(faker.FakeData(proof))
	proof = model.NewProof(proof)
	proof.TenantID = tenantID
	proof.ConnectionID = connectionID
	return proof
}

func fakeEvent(tenantID, connectionID string, jobID *string) *model.Event {
	event := model.NewEvent(nil)
	err2.Check(faker.FakeData(event))

	event = model.NewEvent(event)
	event.TenantID = tenantID
	event.ConnectionID = &connectionID
	event.JobID = utils.CopyStrPtr(jobID)
	return event
}

func fakeJob(
	tenantID, connectionID string,
	protocolConnectionID, protocolCredentialID, protocolProofID, protocolMessageID *string,
) *model.Job {
	job := model.NewJob(nil)
	err2.Check(faker.FakeData(job))
	job = model.NewJob(job)
	job.TenantID = tenantID
	job.ConnectionID = &connectionID
	job.ProtocolConnectionID = utils.CopyStrPtr(protocolConnectionID)
	job.ProtocolCredentialID = utils.CopyStrPtr(protocolCredentialID)
	job.ProtocolProofID = utils.CopyStrPtr(protocolProofID)
	job.ProtocolMessageID = utils.CopyStrPtr(protocolMessageID)
	return job
}

func fakeMessage(tenantID, connectionID string) *model.Message {
	message := model.NewMessage(nil)
	err2.Check(faker.FakeData(message))
	message = model.NewMessage(message)
	message.TenantID = tenantID
	message.ConnectionID = connectionID
	return message
}
