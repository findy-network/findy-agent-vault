package test

import (
	"testing"

	"github.com/findy-network/findy-agent-vault/graph/model"
)

func TestGetMessageConnection(t *testing.T) {
	connection, err := r.BasicMessage().Connection(testContext(), &model.BasicMessage{ID: testMessageID})
	if err != nil {
		t.Errorf("Received unexpected error %s", err)
	}
	if connection == nil {
		t.Errorf("Expecting result, received %v", connection)
	}
}
