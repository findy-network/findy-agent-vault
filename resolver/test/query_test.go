package test

import (
	"context"
	"encoding/base64"
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
	testPaginationErrors(t, "credentials", func(ctx context.Context, after, before *string, first, last *int) error {
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

func TestGetProof(t *testing.T) {
	c, err := r.Query().Proof(testContext(), testProofID)
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

func TestPaginationErrorsGetEvents(t *testing.T) {
	testPaginationErrors(t, "events", func(ctx context.Context, after, before *string, first, last *int) error {
		_, err := r.Query().Events(ctx, after, before, first, last)
		return err
	})
}

func TestResolverGetEvents(t *testing.T) {
	first := 1
	e, err := r.Query().Events(testContext(), nil, nil, &first, nil)
	if err != nil {
		t.Errorf("Received unexpected error %s", err)
	}
	if e == nil || len(e.Edges) == 0 {
		t.Errorf("Expecting result, received %v", e)
	}
}

func TestGetEvent(t *testing.T) {
	e, err := r.Query().Event(testContext(), testEventID)
	if err != nil {
		t.Errorf("Received unexpected error %s", err)
	}
	if e == nil {
		t.Errorf("Expecting result, received %v", e)
	}
}

func TestPaginationErrorsGetJobs(t *testing.T) {
	testPaginationErrors(t, "jobs", func(ctx context.Context, after, before *string, first, last *int) error {
		completed := true
		_, err := r.Query().Jobs(ctx, after, before, first, last, &completed)
		return err
	})
}

func TestResolverGetJobs(t *testing.T) {
	first := 1
	completed := true
	j, err := r.Query().Jobs(testContext(), nil, nil, &first, nil, &completed)
	if err != nil {
		t.Errorf("Received unexpected error %s", err)
	}
	if j == nil || len(j.Edges) == 0 {
		t.Errorf("Expecting result, received %v", j)
	}
}

func TestGetJob(t *testing.T) {
	j, err := r.Query().Job(testContext(), testJobID)
	if err != nil {
		t.Errorf("Received unexpected error %s", err)
	}
	if j == nil {
		t.Errorf("Expecting result, received %v", j)
	}
}

func TestGetUser(t *testing.T) {
	u, err := r.Query().User(testContext())
	if err != nil {
		t.Errorf("Received unexpected error %s", err)
	}
	if u == nil {
		t.Errorf("Expecting result, received %v", u)
	}
}

func TestGetEndpoint(t *testing.T) {
	const expectedLabel = "findy-issuer"
	// plain json string
	e, err := r.Query().Endpoint(testContext(), testInvitation)
	if err != nil {
		t.Errorf("Received unexpected error %s", err)
	}
	if e == nil || e.Label != expectedLabel {
		t.Errorf("Expecting valid result, received %v", e)
	}

	// base64 encoded string
	e, err = r.Query().Endpoint(testContext(), base64.StdEncoding.EncodeToString([]byte(testInvitation)))
	if err != nil {
		t.Errorf("Received unexpected error %s", err)
	}
	if e == nil || e.Label != expectedLabel {
		t.Errorf("Expecting valid result, received %v", e)
	}
}
