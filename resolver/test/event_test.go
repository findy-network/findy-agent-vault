package test

import (
	"testing"

	"github.com/findy-network/findy-agent-vault/graph/model"
)

func TestGetEventConnection(t *testing.T) {
	beforeEach(t)

	connection, err := r.Event().Connection(testContext(), &model.Event{ID: testEventID})
	if err != nil {
		t.Errorf("Received unexpected error %s", err)
	}
	if connection == nil {
		t.Errorf("Expecting result, received %v", connection)
	}
}

func TestGetEventJob(t *testing.T) {
	beforeEach(t)

	job, err := r.Event().Job(testContext(), &model.Event{ID: testEventID})
	if err != nil {
		t.Errorf("Received unexpected error %s", err)
	}
	if job == nil {
		t.Errorf("Expecting result, received %v", job)
	}
}
