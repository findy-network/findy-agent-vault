package resolver

import (
	"context"
	"testing"
)

func TestPaginationErrorsGetConnections(t *testing.T) {
	testPaginationErrors(t, "connections", func(ctx context.Context, after, before *string, first, last *int) error {
		_, err := r.Query().Connections(ctx, after, before, first, last)
		return err
	})
}

func TestResolverGetConnections(t *testing.T) {
	first := 1
	c, err := r.Query().Connections(testContext(), nil, nil, &first, nil)
	if err != nil {
		t.Errorf("Received unexpected error %s", err)
	}
	if c == nil || len(c.Edges) == 0 {
		t.Errorf("Expecting result, received %v", c)
	}
}

func TestGetConnection(t *testing.T) {
	c, err := r.Query().Connection(testContext(), testConnectionID)
	if err != nil {
		t.Errorf("Received unexpected error %s", err)
	}
	if c == nil {
		t.Errorf("Expecting result, received %v", c)
	}
}

func TestPaginationErrorsGetCredentials(t *testing.T) {
	testPaginationErrors(t, "connections", func(ctx context.Context, after, before *string, first, last *int) error {
		_, err := r.Query().Credentials(ctx, after, before, first, last)
		return err
	})
}

func TestResolverGetCredentials(t *testing.T) {
	first := 1
	c, err := r.Query().Credentials(testContext(), nil, nil, &first, nil)
	if err != nil {
		t.Errorf("Received unexpected error %s", err)
	}
	if c == nil || len(c.Edges) == 0 {
		t.Errorf("Expecting result, received %v", c)
	}
}

func TestGetCredential(t *testing.T) {
	c, err := r.Query().Credential(testContext(), testCredentialID)
	if err != nil {
		t.Errorf("Received unexpected error %s", err)
	}
	if c == nil {
		t.Errorf("Expecting result, received %v", c)
	}
}

func TestPaginationErrorsGetEvents(t *testing.T) {
	testPaginationErrors(t, "connections", func(ctx context.Context, after, before *string, first, last *int) error {
		_, err := r.Query().Events(ctx, after, before, first, last)
		return err
	})
}

func TestResolverGetEvents(t *testing.T) {
	first := 1
	c, err := r.Query().Events(testContext(), nil, nil, &first, nil)
	if err != nil {
		t.Errorf("Received unexpected error %s", err)
	}
	if c == nil || len(c.Edges) == 0 {
		t.Errorf("Expecting result, received %v", c)
	}
}

func TestGetEvent(t *testing.T) {
	c, err := r.Query().Event(testContext(), testEventID)
	if err != nil {
		t.Errorf("Received unexpected error %s", err)
	}
	if c == nil {
		t.Errorf("Expecting result, received %v", c)
	}
}

func TestGetMessage(t *testing.T) {
	c, err := r.Query().Message(testContext(), testMessageID)
	if err != nil {
		t.Errorf("Received unexpected error %s", err)
	}
	if c == nil {
		t.Errorf("Expecting result, received %v", c)
	}
}
