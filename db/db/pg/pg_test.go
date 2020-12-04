package pg

import (
	"os"
	"testing"
	"time"

	"github.com/findy-network/findy-agent-vault/db/db"
	"github.com/findy-network/findy-agent-vault/db/model"
	"github.com/findy-network/findy-agent-vault/tools/utils"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var pgDb db.Db
var m *migrate.Migrate

func setup() {
	utils.SetLogDefaults()
	pgDb = InitDb("file://../../migrations", "5433", true)
}

func teardown() {
	pgDb.Close()
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func TestAddAgent(t *testing.T) {
	var (
		agentID = "agentID"
		label   = "agentLabel"
	)

	// Add data
	if err := pgDb.AddAgent(agentID, label); err != nil {
		t.Errorf("Failed to add agent %s", err.Error())
	}

	// Error with duplicate id
	err := pgDb.AddAgent(agentID, label)
	if err == nil {
		t.Errorf("Expecting duplicate key error")
	}

	if pgErr, ok := err.(*PgError); ok {
		if pgErr.code != PgErrorUniqueViolation {
			t.Errorf("Expecting duplicate key error %s", pgErr.code)
		}
	} else {
		t.Errorf("Expecting pg error %v", err)
	}

	var validateAgent = func(a *model.Agent) {
		if a == nil {
			t.Errorf("Expecting result, agent is nil")
			return
		}
		if a.AgentID != agentID {
			t.Errorf("Agent id mismatch expected %s got %s", agentID, a.AgentID)
		}
		if a.Label != label {
			t.Errorf("Agent label mismatch expected %s got %s", label, a.Label)
		}
		if a.ID == "" {
			t.Errorf("Invalid agent id %s", a.ID)
		}
		if time.Since(a.Created) > time.Second {
			t.Errorf("Timestamp not in threshold %v", a.Created)
		}
	}

	// Get data for agent id
	a1, err := pgDb.GetAgent(nil, &agentID)
	if err != nil {
		t.Errorf("Error fetching agent %s", err.Error())
	} else {
		validateAgent(a1)
	}

	// Get data for id
	a2, err := pgDb.GetAgent(&a1.ID, nil)
	if err != nil {
		t.Errorf("Error fetching agent %s", err.Error())
	} else {
		validateAgent(a2)
	}
}
