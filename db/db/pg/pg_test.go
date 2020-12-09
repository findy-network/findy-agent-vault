package pg

import (
	"os"
	"testing"

	"github.com/findy-network/findy-agent-vault/db/db"
	"github.com/findy-network/findy-agent-vault/db/model"
	"github.com/findy-network/findy-agent-vault/utils"
)

var (
	pgDB             db.DB
	testTenantID     string
	testAgentID      string
	testConnectionID string
	testConnection   *model.Connection = &model.Connection{
		OurDid:        "ourDid",
		TheirDid:      "theirDid",
		TheirEndpoint: "theirEndpoint",
		TheirLabel:    "theirLabel",
		Invited:       false,
	}
)

func setup() {
	utils.SetLogDefaults()
	pgDB = InitDB("file://../../migrations", "5433", true)

	testAgent := &model.Agent{AgentID: "testAgentID", Label: "testAgent"}

	a, err := pgDB.AddAgent(testAgent)
	if err != nil {
		panic(err)
	}
	testTenantID = a.ID
	testAgentID = a.AgentID
	testConnection.TenantID = testTenantID

	c, err := pgDB.AddConnection(testConnection)
	if err != nil {
		panic(err)
	}
	testConnectionID = c.ID
}

func teardown() {
	pgDB.Close()
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}
