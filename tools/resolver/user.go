package resolver

import (
	"context"

	"github.com/findy-network/findy-agent-api/tools/data"

	"github.com/findy-network/findy-agent-api/graph/model"
)

func (r *queryResolver) User(_ context.Context) (*model.User, error) {
	return data.State.User.ToNode(), nil
}
