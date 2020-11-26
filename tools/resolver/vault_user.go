package resolver

import (
	"context"

	"github.com/findy-network/findy-agent-vault/graph/model"
)

func (r *queryResolver) User(ctx context.Context) (*model.User, error) {
	return state.User.ToNode(), nil
}
