package pg

import (
	"os"
	"testing"

	"github.com/findy-network/findy-agent-vault/db/db"
	"github.com/findy-network/findy-agent-vault/db/model"
	"github.com/findy-network/findy-agent-vault/tools/utils"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var (
	pgDB         db.Db
	testTenantId string
)

func setup() {
	utils.SetLogDefaults()
	pgDB = InitDB("file://../../migrations", "5433", true)

	testAgent := &model.Agent{AgentID: "testAgentID", Label: "testAgent"}

	a, err := pgDB.AddAgent(testAgent)
	if err != nil {
		panic(err)
	}
	testTenantId = a.ID
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
