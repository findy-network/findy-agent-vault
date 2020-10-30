package agency

import (
	"os"
	"testing"
)

type agencyListener struct{}

func (l *agencyListener) AddConnection(id, ourDID, theirDID, theirEndpoint, theirLabel string, initiatedByUs bool) {

}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func setup() {
	Instance.Init(&agencyListener{})
}

func teardown() {
}
