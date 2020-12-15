package resolver

import (
	"os"
	"testing"

	"github.com/findy-network/findy-agent-vault/db/fake"

	"github.com/findy-network/findy-agent-vault/db/store/test"
	"github.com/findy-network/findy-agent-vault/utils"
)

var (
	r                *Resolver
	testConnectionID string
	testCredentialID string
)

func setup() {
	utils.SetLogDefaults()
	r = InitResolver(true)
	size := 5
	a, c := test.AddAgentAndConnections(r.db, fake.FakeCloudDID, size)
	testConnectionID = c[0].ID

	cr := fake.AddCredentials(r.db, a.ID, c[0].ID, size)
	testCredentialID = cr[0].ID

}

func teardown() {
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}
