package test

import (
	"testing"

	"github.com/findy-network/findy-agent-vault/graph/model"
)

func TestGetJobOutput(t *testing.T) {
	output, err := r.Job().Output(
		testContext(),
		&model.Job{ID: testJobID, Protocol: model.ProtocolTypeCredential},
	)
	if err != nil {
		t.Errorf("Received unexpected error %s", err)
	}
	if output == nil || output.Credential == nil {
		t.Errorf("Expecting result, received %v", output)
	}
}
