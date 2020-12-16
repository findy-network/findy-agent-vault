package resolver

import (
	"testing"

	"github.com/findy-network/findy-agent-vault/graph/model"
)

func TestGetEventConnection(t *testing.T) {
	connection, err := r.Event().Connection(testContext(), &model.Event{ID: testEventID})
	if err != nil {
		t.Errorf("Received unexpected error %s", err)
	}
	if connection == nil {
		t.Errorf("Expecting result, received %v", connection)
	}
}
