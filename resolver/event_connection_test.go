package resolver

import (
	"testing"

	"github.com/findy-network/findy-agent-vault/graph/model"
)

func TestEventConnectionTotalCount(t *testing.T) {
	c, err := r.EventConnection().TotalCount(testContext(), &model.EventConnection{ConnectionID: nil})
	if err != nil {
		t.Errorf("Received unexpected error %s", err)
	}
	if c != totalCount {
		t.Errorf("Mismatch in count exp: %d, got: %d", totalCount, c)
	}
}

func TestEventConnectionTotalCountForConnection(t *testing.T) {
	c, err := r.EventConnection().TotalCount(testContext(), &model.EventConnection{ConnectionID: &testConnectionID})
	if err != nil {
		t.Errorf("Received unexpected error %s", err)
	}
	if c != totalCount {
		t.Errorf("Mismatch in count exp: %d, got: %d", totalCount, c)
	}
}
