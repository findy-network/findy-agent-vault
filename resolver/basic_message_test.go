package resolver

import (
	"testing"

	"github.com/findy-network/findy-agent-vault/graph/model"
)

func TestGetMessageConnection(t *testing.T) {
	j, err := r.BasicMessage().Connection(testContext(), &model.BasicMessage{ID: testMessageID})
	if err != nil {
		t.Errorf("Received unexpected error %s", err)
	}
	if j == nil {
		t.Errorf("Expecting result, received %v", j)
	}
}
