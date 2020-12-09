package fake

import (
	crand "crypto/rand"
	"math/big"
	"reflect"

	"github.com/bxcodec/faker/v3"
	"github.com/findy-network/findy-agent-vault/db/db"
	"github.com/findy-network/findy-agent-vault/db/model"
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

func addProviders(tenantID string) {
	_ = faker.AddProvider("organisationLabel", func(v reflect.Value) (interface{}, error) {
		orgs := []string{"Bank", "Ltd", "Agency", "Company", "United"}
		index := random(len(orgs))
		return faker.LastName() + " " + orgs[index], nil
	})

	_ = faker.AddProvider("tenantId", func(v reflect.Value) (interface{}, error) {
		return tenantID, nil
	})
}

func AddConnections(store db.DB, tenantID string, count int) []*model.Connection {
	addProviders(tenantID)

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
