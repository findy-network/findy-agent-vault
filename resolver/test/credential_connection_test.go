package test

import (
	"testing"

	"github.com/findy-network/findy-agent-vault/graph/model"
)

func TestCredentialConnectionTotalCount(t *testing.T) {
	const user = "TestCredentialConnectionTotalCount"
	beforeEachWithID(t, user)

	c, err := r.CredentialConnection().TotalCount(testContextForUser(user), &model.CredentialConnection{ConnectionID: nil})
	if err != nil {
		t.Errorf("Received unexpected error %s", err)
	}
	if c != totalCount {
		t.Errorf("Mismatch in count exp: %d, got: %d", totalCount, c)
	}
}

func TestCredentialConnectionTotalCountForConnection(t *testing.T) {
	beforeEach(t)

	c, err := r.CredentialConnection().TotalCount(testContext(), &model.CredentialConnection{ConnectionID: &testConnectionID})
	if err != nil {
		t.Errorf("Received unexpected error %s", err)
	}
	if c != totalCount {
		t.Errorf("Mismatch in count exp: %d, got: %d", totalCount, c)
	}
}
