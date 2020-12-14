package test

import (
	"os"
	"testing"

	"github.com/findy-network/findy-agent-vault/db/model"
	"github.com/findy-network/findy-agent-vault/db/store"
	"github.com/findy-network/findy-agent-vault/db/store/mock"
	"github.com/findy-network/findy-agent-vault/db/store/pg"
	graph "github.com/findy-network/findy-agent-vault/graph/model"
	"github.com/findy-network/findy-agent-vault/utils"
)

type testableDB struct {
	db               store.DB
	name             string
	testTenantID     string
	testAgentID      string
	testConnectionID string
	testConnection   *model.Connection
}

var (
	DBs            []*testableDB
	testCredential *model.Credential = &model.Credential{
		Role:          graph.CredentialRoleHolder,
		SchemaID:      "schemaId",
		CredDefID:     "credDefId",
		InitiatedByUs: false,
		Attributes: []*graph.CredentialValue{
			{Name: "name1", Value: "value1"},
			{Name: "name2", Value: "value2"},
		},
	}
	testEvent *model.Event = &model.Event{
		Description: "event desc",
		Read:        false,
	}
)

func setup() {
	utils.SetLogDefaults()

	testAgent := model.NewAgent()
	testAgent.AgentID = "testAgentID"
	testAgent.Label = "testAgent"

	testConnection := model.NewConnection(nil)
	testConnection.OurDid = "ourDid"
	testConnection.TheirDid = "theirDid"
	testConnection.TheirEndpoint = "theirEndpoint"
	testConnection.TheirLabel = "theirLabel"
	testConnection.Invited = false

	DBs = append(DBs, []*testableDB{{
		db:             pg.InitDB("file://../../migrations", "5433", true),
		name:           "pg",
		testConnection: testConnection,
	}, {
		db:             mock.InitState(),
		name:           "mock",
		testConnection: testConnection,
	},
	}...)

	for index := range DBs {
		s := DBs[index]

		a, err := s.db.AddAgent(testAgent)
		if err != nil {
			panic(err)
		}
		s.testTenantID = a.ID
		s.testAgentID = a.AgentID
		s.testConnection.TenantID = s.testTenantID

		c, err := s.db.AddConnection(testConnection)
		if err != nil {
			panic(err)
		}
		s.testConnectionID = c.ID
		s.testConnection = c
	}
}

func teardown() {
	for _, s := range DBs {
		s.db.Close()
	}
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func addAgentAndConnections(agentID string, s *testableDB) (*model.Agent, []*model.Connection) {
	// add new agent with no pre-existing event s
	ctAgent := model.NewAgent()
	ctAgent.AgentID = agentID
	ctAgent.Label = "testAgent"
	a, err := s.db.AddAgent(ctAgent)
	if err != nil {
		panic(err)
	}
	// add new connections
	connCount := 3
	connections := make([]*model.Connection, connCount)
	for i := 0; i < connCount; i++ {
		c, err := s.db.AddConnection(s.testConnection)
		if err != nil {
			panic(err)
		}
		connections[i] = c
	}
	return a, connections
}
