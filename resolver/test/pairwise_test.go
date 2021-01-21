package test

import (
	"context"
	"testing"

	"github.com/findy-network/findy-agent-vault/graph/model"
)

func TestPaginationErrorsGetConnectionCredentials(t *testing.T) {
	testPaginationErrors(t, "connection credentials", func(ctx context.Context, after, before *string, first, last *int) error {
		_, err := r.Pairwise().Credentials(ctx, &model.Pairwise{ID: testConnectionID}, after, before, first, last)
		return err
	})
}

func TestResolverGetConnectionCredentials(t *testing.T) {
	first := 1
	c, err := r.Pairwise().Credentials(testContext(), &model.Pairwise{ID: testConnectionID}, nil, nil, &first, nil)
	if err != nil {
		t.Errorf("Received unexpected error %s", err)
	}
	if c == nil || len(c.Edges) == 0 {
		t.Errorf("Expecting result, received %v", c)
	}
}

func TestPaginationErrorsGetConnectionProofs(t *testing.T) {
	testPaginationErrors(t, "connection proofs", func(ctx context.Context, after, before *string, first, last *int) error {
		_, err := r.Pairwise().Proofs(ctx, &model.Pairwise{ID: testConnectionID}, after, before, first, last)
		return err
	})
}

func TestResolverGetConnectionProofs(t *testing.T) {
	first := 1
	c, err := r.Pairwise().Proofs(testContext(), &model.Pairwise{ID: testConnectionID}, nil, nil, &first, nil)
	if err != nil {
		t.Errorf("Received unexpected error %s", err)
	}
	if c == nil || len(c.Edges) == 0 {
		t.Errorf("Expecting result, received %v", c)
	}
}

func TestPaginationErrorsGetConnectionMessages(t *testing.T) {
	testPaginationErrors(t, "connection messages", func(ctx context.Context, after, before *string, first, last *int) error {
		_, err := r.Pairwise().Messages(ctx, &model.Pairwise{ID: testConnectionID}, after, before, first, last)
		return err
	})
}

func TestResolverGetConnectionMessages(t *testing.T) {
	first := 1
	c, err := r.Pairwise().Messages(testContext(), &model.Pairwise{ID: testConnectionID}, nil, nil, &first, nil)
	if err != nil {
		t.Errorf("Received unexpected error %s", err)
	}
	if c == nil || len(c.Edges) == 0 {
		t.Errorf("Expecting result, received %v", c)
	}
}

func TestPaginationErrorsGetConnectionEvents(t *testing.T) {
	testPaginationErrors(t, "connection events", func(ctx context.Context, after, before *string, first, last *int) error {
		_, err := r.Pairwise().Events(ctx, &model.Pairwise{ID: testConnectionID}, after, before, first, last)
		return err
	})
}

func TestResolverGetConnectionEvents(t *testing.T) {
	first := 1
	c, err := r.Pairwise().Events(testContext(), &model.Pairwise{ID: testConnectionID}, nil, nil, &first, nil)
	if err != nil {
		t.Errorf("Received unexpected error %s", err)
	}
	if c == nil || len(c.Edges) == 0 {
		t.Errorf("Expecting result, received %v", c)
	}
}

func TestPaginationErrorsGetConnectionJobs(t *testing.T) {
	testPaginationErrors(t, "connection jobs", func(ctx context.Context, after, before *string, first, last *int) error {
		completed := true
		_, err := r.Pairwise().Jobs(ctx, &model.Pairwise{ID: testConnectionID}, after, before, first, last, &completed)
		return err
	})
}

func TestResolverGetConnectionJobs(t *testing.T) {
	first := 1
	completed := true
	j, err := r.Pairwise().Jobs(testContext(), &model.Pairwise{ID: testConnectionID}, nil, nil, &first, nil, &completed)
	if err != nil {
		t.Errorf("Received unexpected error %s", err)
	}
	if j == nil || len(j.Edges) == 0 {
		t.Errorf("Expecting result, received %v", j)
	}
}
