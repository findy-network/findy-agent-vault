package resolver

import (
	"context"
	"os"
	"testing"

	"github.com/findy-network/findy-agent-vault/db/fake"
	"github.com/findy-network/findy-agent-vault/server"

	"github.com/findy-network/findy-agent-vault/db/store/test"
	"github.com/findy-network/findy-agent-vault/utils"
)

var (
	r                *Resolver
	testConnectionID string
	testCredentialID string
	testEventID      string
	totalCount       = 5
)

func testContext() context.Context {
	u := server.CreateTestToken("test")
	ctx := context.WithValue(context.Background(), "user", u)
	return ctx
}

func setup() {
	utils.SetLogDefaults()
	r = InitResolver(true)
	size := totalCount
	a, c := test.AddAgentAndConnections(r.db, fake.FakeCloudDID, size)
	testConnectionID = c[0].ID

	cr := fake.AddCredentials(r.db, a.ID, c[0].ID, size)
	testCredentialID = cr[0].ID

	ev := fake.AddEvents(r.db, a.ID, c[0].ID, size)
	testEventID = ev[0].ID
}

func teardown() {
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}
