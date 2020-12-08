package fake

import (
	crand "crypto/rand"
	"math/big"
	"reflect"

	"github.com/bxcodec/faker/v3"
	"github.com/findy-network/findy-agent-vault/db/db"
	"github.com/findy-network/findy-agent-vault/db/model"
	"github.com/findy-network/findy-agent-vault/utils"
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

func AddConnections(db db.Db, tenantID string, count int) []*model.Connection {
	addProviders(tenantID)

	connections := make([]*model.Connection, count)
	for i := 0; i < count; i++ {
		connection := &model.Connection{}
		err := faker.FakeData(connection)
		if err != nil {
			panic(err)
		}
		connections[i] = connection
	}

	newConnections := make([]*model.Connection, count)
	for index, connection := range connections {
		c, err := db.AddConnection(connection)
		if err != nil {
			panic(err)
		}
		newConnections[index] = c
	}

	utils.LogMed().Infof("Generated %d connections for tenant %s", len(newConnections), tenantID)

	return newConnections
}

func AddAgent(db db.Db) *model.Agent {
	_ = faker.AddProvider("agentId", func(v reflect.Value) (interface{}, error) {
		return FakeCloudDID, nil
	})
	var err error
	agent := &model.Agent{}
	faker.FakeData(&agent)
	if agent, err = db.AddAgent(agent); err != nil {
		panic(err)
	}
	utils.LogMed().Infof("Generated tenant %s with agent id %s", agent.ID, agent.AgentID)
	return agent
}

func AddData(db db.Db) {
	agent := AddAgent(db)

	AddConnections(db, agent.ID, 5)
}
