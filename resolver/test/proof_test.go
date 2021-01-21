package test

import (
	"testing"

	"github.com/findy-network/findy-agent-vault/graph/model"
)

func TestGetProofConnection(t *testing.T) {
	connection, err := r.Proof().Connection(testContext(), &model.Proof{ID: testProofID})
	if err != nil {
		t.Errorf("Received unexpected error %s", err)
	}
	if connection == nil {
		t.Errorf("Expecting result, received %v", connection)
	}
}
