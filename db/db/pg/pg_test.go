package pg

import (
	"math"
	"os"
	"testing"
	"time"

	"github.com/findy-network/findy-agent-vault/db/db"
	"github.com/findy-network/findy-agent-vault/db/model"
	"github.com/findy-network/findy-agent-vault/utils"
)

var (
	pgDB         db.Db
	testTenantID string
	testAgentID  string
)

func ceilTimestamp(ts *time.Time) uint64 {
	return uint64(math.Ceil(float64(ts.UnixNano()) / float64(time.Second.Nanoseconds())))
}

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
