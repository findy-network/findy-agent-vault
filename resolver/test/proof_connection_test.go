package test

import (
	"testing"

	"github.com/findy-network/findy-agent-vault/graph/model"
)

func TestProofConnectionTotalCountForConnection(t *testing.T) {
	c, err := r.ProofConnection().TotalCount(testContext(), &model.ProofConnection{ConnectionID: &testConnectionID})
	if err != nil {
		t.Errorf("Received unexpected error %s", err)
	}
	if c != totalCount {
		t.Errorf("Mismatch in count exp: %d, got: %d", totalCount, c)
	}
}
