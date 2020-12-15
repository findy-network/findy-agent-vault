package resolver

import (
	"testing"

	"github.com/findy-network/findy-agent-vault/graph/model"
)

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
