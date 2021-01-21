package test

import (
	"testing"

	"github.com/findy-network/findy-agent-vault/graph/model"
)

func TestCredentialConnectionTotalCount(t *testing.T) {
	c, err := r.CredentialConnection().TotalCount(testContext(), &model.CredentialConnection{ConnectionID: nil})
	if err != nil {
		t.Errorf("Received unexpected error %s", err)
	}
	if c != totalCount {
		t.Errorf("Mismatch in count exp: %d, got: %d", totalCount, c)
	}
}

func TestCredentialConnectionTotalCountForConnection(t *testing.T) {
	c, err := r.CredentialConnection().TotalCount(testContext(), &model.CredentialConnection{ConnectionID: &testConnectionID})
	if err != nil {
		t.Errorf("Received unexpected error %s", err)
	}
	if c != totalCount {
		t.Errorf("Mismatch in count exp: %d, got: %d", totalCount, c)
	}
}
