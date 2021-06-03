package test

import (
	"testing"

	"github.com/findy-network/findy-agent-vault/graph/model"
)

func TestEventConnectionTotalCount(t *testing.T) {
	const user = "TestEventConnectionTotalCount"
	beforeEachWithID(t, user)

	c, err := r.EventConnection().TotalCount(testContextForUser(user), &model.EventConnection{ConnectionID: nil})
	if err != nil {
		t.Errorf("Received unexpected error %s", err)
	}
	if c != totalCount {
		t.Errorf("Mismatch in count exp: %d, got: %d", totalCount, c)
	}
}

func TestEventConnectionTotalCountForConnection(t *testing.T) {
	beforeEach(t)

	c, err := r.EventConnection().TotalCount(testContext(), &model.EventConnection{ConnectionID: &testConnectionID})
	if err != nil {
		t.Errorf("Received unexpected error %s", err)
	}
	if c != totalCount {
		t.Errorf("Mismatch in count exp: %d, got: %d", totalCount, c)
	}
}
