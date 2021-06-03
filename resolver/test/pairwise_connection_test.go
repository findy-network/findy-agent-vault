package test

import (
	"testing"

	"github.com/findy-network/findy-agent-vault/graph/model"
)

func TestPairwiseConnectionTotalCount(t *testing.T) {
	const user = "TestPairwiseConnectionTotalCount"
	beforeEachWithID(t, user)

	c, err := r.PairwiseConnection().TotalCount(testContextForUser(user), &model.PairwiseConnection{})
	if err != nil {
		t.Errorf("Received unexpected error %s", err)
	}
	if c != totalCount {
		t.Errorf("Mismatch in count exp: %d, got: %d", totalCount, c)
	}
}
