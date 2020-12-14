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

func AddEvents(db store.DB, tenantID, connectionID string, count int) []*model.Event {
	events := make([]*model.Event, count)
	for i := 0; i < count; i++ {
		event := fakeEvent(tenantID, connectionID)
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

func fakeAgent() *model.Agent {
	agent := model.NewAgent()
	err2.Check(faker.FakeData(&agent))
	return agent.Copy()
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

func fakeEvent(tenantID, connectionID string) *model.Event {
	event := model.NewEvent(nil)
	err2.Check(faker.FakeData(event))
	event = model.NewEvent(event)
	event.TenantID = tenantID
	event.ConnectionID = &connectionID
	return event
}
