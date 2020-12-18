package resolver

import (
	"testing"

	"github.com/findy-network/findy-agent-vault/graph/model"
)

func TestMarkEventRead(t *testing.T) {
	event, err := r.Mutation().MarkEventRead(testContext(), model.MarkReadInput{ID: testEventID})
	if err != nil {
		t.Errorf("Received unexpected error %s", err)
	}
	if event == nil {
		t.Errorf("Expecting result, received %v", event)
	}
}
