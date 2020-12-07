package fake

import (
	crand "crypto/rand"
	"math/big"
	"reflect"

	"github.com/bxcodec/faker/v3"
	"github.com/findy-network/findy-agent-vault/db/db"
	"github.com/findy-network/findy-agent-vault/db/model"
)

func random(n int) int {
	val, err := crand.Int(crand.Reader, big.NewInt(int64(n)))
	if err != nil {
		panic(err)
	}
	return int(val.Int64())
}

func addProviders(agentID string) {
	_ = faker.AddProvider("organisationLabel", func(v reflect.Value) (interface{}, error) {
		orgs := []string{"Bank", "Ltd", "Agency", "Company", "United"}
		index := random(len(orgs))
		return faker.LastName() + " " + orgs[index], nil
	})

	_ = faker.AddProvider("tenantId", func(v reflect.Value) (interface{}, error) {
		return agentID, nil
	})
}

func AddConnections(db db.Db, agentID string, count int) []*model.Connection {
	addProviders(agentID)

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

	return newConnections
}

func AddData(db db.Db) {
	var err error

	agent := &model.Agent{}
	faker.FakeData(&agent)
	if agent, err = db.AddAgent(agent); err != nil {
		panic(err)
	}

	AddConnections(db, agent.ID, 5)
}
