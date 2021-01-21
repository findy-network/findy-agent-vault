package test

import (
	"testing"

	"github.com/findy-network/findy-agent-vault/graph/model"
)

func TestGetCredentialConnection(t *testing.T) {
	connection, err := r.Credential().Connection(testContext(), &model.Credential{ID: testCredentialID})
	if err != nil {
		t.Errorf("Received unexpected error %s", err)
	}
	if connection == nil {
		t.Errorf("Expecting result, received %v", connection)
	}
}
