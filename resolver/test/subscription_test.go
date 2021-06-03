package test

import (
	"testing"
)

func TestSubscribeEventAdded(t *testing.T) {
	beforeEach(t)

	channel, err := r.Subscription().EventAdded(testContext())
	if err != nil {
		t.Errorf("Received unexpected error %s", err)
	}
	if channel == nil {
		t.Errorf("Expecting result, received %v", channel)
	}
}

// TODO: add test for sending events
