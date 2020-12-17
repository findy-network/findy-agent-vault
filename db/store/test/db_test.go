package test

import (
	"os"
	"testing"
	"time"

	"github.com/findy-network/findy-agent-vault/db/model"
	"github.com/findy-network/findy-agent-vault/db/store"
	"github.com/findy-network/findy-agent-vault/db/store/mock"
	"github.com/findy-network/findy-agent-vault/db/store/pg"
	graph "github.com/findy-network/findy-agent-vault/graph/model"
	"github.com/findy-network/findy-agent-vault/utils"
)

const testAgentLabel = "testAgent"

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
	testCredential *model.Credential = model.NewCredential(&model.Credential{
		Role:          graph.CredentialRoleHolder,
		SchemaID:      "schemaId",
		CredDefID:     "credDefId",
		InitiatedByUs: false,
		Attributes: []*graph.CredentialValue{
			{Name: "name1", Value: "value1"},
			{Name: "name2", Value: "value2"},
		},
	})
	testProof *model.Proof = model.NewProof(&model.Proof{
		Role:          graph.ProofRoleProver,
		InitiatedByUs: false,
		Result:        true,
		Attributes: []*graph.ProofAttribute{
			{Name: "name1", Value: nil, CredDefID: "cred_def_id"},
			{Name: "name2", Value: nil, CredDefID: "cred_def_id"},
		},
	})
	testMessage *model.Message = model.NewMessage(&model.Message{
		Message:   "msg content",
		SentByMe:  false,
		Delivered: nil,
	})
	testEvent *model.Event = model.NewEvent(&model.Event{
		Description: "event desc",
		Read:        false,
	})
	testJob *model.Job = model.NewJob("", "", &model.Job{
		ProtocolType: graph.ProtocolTypeConnection,
		Status:       graph.JobStatusWaiting,
		Result:       graph.JobResultNone,
	})
)

func validateTimestap(t *testing.T, exp, got *time.Time, name string) {
	fail := false
	if got != exp {
		fail = true
		if got != nil && exp != nil && got.Sub(*exp) < time.Microsecond {
			fail = false
		}
	}
	if fail {
		t.Errorf("Object %s mismatch expected %s got %s", name, exp, got)
	}
}

func validateStrPtr(t *testing.T, exp, got *string, name string) {
	fail := false
	if got != exp {
		fail = true
		if got != nil && exp != nil {
			if *got != *exp {
				t.Errorf("Object %s mismatch expected %s got %s", name, *exp, *got)
			} else {
				fail = false
			}
		}
	}
	if fail {
		t.Errorf("Object %s mismatch expected %v got %v", name, exp, got)
	}
}

func setup() {
	utils.SetLogDefaults()

	testAgent := model.NewAgent(nil)
	testAgent.AgentID = "testAgentID"
	testAgent.Label = testAgentLabel

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
