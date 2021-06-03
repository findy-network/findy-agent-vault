package test

import (
	"testing"

	"github.com/findy-network/findy-agent-vault/graph/model"
)

func TestMessageConnectionTotalCountForConnection(t *testing.T) {
	beforeEach(t)

	c, err := r.BasicMessageConnection().TotalCount(testContext(), &model.BasicMessageConnection{ConnectionID: &testConnectionID})
	if err != nil {
		t.Errorf("Received unexpected error %s", err)
	}
	if c != totalCount {
		t.Errorf("Mismatch in count exp: %d, got: %d", totalCount, c)
	}
}
