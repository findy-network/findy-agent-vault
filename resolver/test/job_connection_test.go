package test

import (
	"testing"

	"github.com/findy-network/findy-agent-vault/graph/model"
)

func TestJobConnectionTotalCount(t *testing.T) {
	const user = "TestJobConnectionTotalCount"
	beforeEachWithID(t, user)

	completed := true
	j, err := r.JobConnection().TotalCount(testContextForUser(user), &model.JobConnection{ConnectionID: nil, Completed: &completed})
	if err != nil {
		t.Errorf("Received unexpected error %s", err)
	}
	if j != totalCount {
		t.Errorf("Mismatch in count exp: %d, got: %d", totalCount, j)
	}
}

func TestJobConnectionTotalCountForConnection(t *testing.T) {
	beforeEach(t)

	completed := true
	j, err := r.JobConnection().TotalCount(testContext(), &model.JobConnection{ConnectionID: &testConnectionID, Completed: &completed})
	if err != nil {
		t.Errorf("Received unexpected error %s", err)
	}
	if j != totalCount {
		t.Errorf("Mismatch in count exp: %d, got: %d", totalCount, j)
	}
}
