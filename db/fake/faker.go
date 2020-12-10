package fake

import (
	crand "crypto/rand"
	"math/big"
	"reflect"
	"time"

	"github.com/bxcodec/faker/v3"
	"github.com/findy-network/findy-agent-vault/db/db"
	"github.com/findy-network/findy-agent-vault/db/model"
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

func AddCredentials(store db.DB, tenantID, connectionID string, count int) []*model.Credential {
	_ = faker.AddProvider("nil", func(v reflect.Value) (interface{}, error) {
		return nil, nil
	})
	_ = faker.AddProvider("credentialAttributes", func(v reflect.Value) (interface{}, error) {
		return []*graph.CredentialValue{
			{Name: "name1", Value: "value1"},
			{Name: "name2", Value: "value2"},
			{Name: "name3", Value: "value3"},
		}, nil
	})

	credentials := make([]*model.Credential, count)
	for i := 0; i < count; i++ {
		credential := &model.Credential{}
		err2.Check(faker.FakeData(credential))
		credential.TenantID = tenantID
		credential.ConnectionID = connectionID
		credentials[i] = credential
	}

	newCredentials := make([]*model.Credential, count)
	for index, credential := range credentials {
		c, err := store.AddCredential(credential)
		err2.Check(err)

		now := time.Now().UTC()
		c.Approved = &now
		c.Issued = &now
		_, err = store.UpdateCredential(c)
		err2.Check(err)

		newCredentials[index] = c
	}

	utils.LogMed().Infof("Generated %d credentials for tenant %s", len(newCredentials), tenantID)

	return newCredentials
}

func AddConnections(store db.DB, tenantID string, count int) []*model.Connection {
	_ = faker.AddProvider("organisationLabel", func(v reflect.Value) (interface{}, error) {
		orgs := []string{"Bank", "Ltd", "Agency", "Company", "United"}
		index := random(len(orgs))
		return faker.LastName() + " " + orgs[index], nil
	})

	_ = faker.AddProvider("tenantId", func(v reflect.Value) (interface{}, error) {
		return tenantID, nil
	})

	connections := make([]*model.Connection, count)
	for i := 0; i < count; i++ {
		connection := &model.Connection{}
		err2.Check(faker.FakeData(connection))
		connections[i] = connection
	}

	newConnections := make([]*model.Connection, count)
	for index, connection := range connections {
		c, err := store.AddConnection(connection)
		err2.Check(err)
		newConnections[index] = c
	}

	utils.LogMed().Infof("Generated %d connections for tenant %s", len(newConnections), tenantID)

	return newConnections
}

func AddAgent(store db.DB) *model.Agent {
	_ = faker.AddProvider("agentId", func(v reflect.Value) (interface{}, error) {
		return FakeCloudDID, nil
	})
	var err error
	agent := &model.Agent{}
	err2.Check(faker.FakeData(&agent))

	agent, err = store.AddAgent(agent)
	err2.Check(err)

	utils.LogMed().Infof("Generated tenant %s with agent id %s", agent.ID, agent.AgentID)
	return agent
}

func AddData(store db.DB) {
	agent := AddAgent(store)

	AddConnections(store, agent.ID, 5)
}
