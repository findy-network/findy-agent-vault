package resolver

import (
	"context"

	"github.com/findy-network/findy-agent-vault/tools/data"

	"github.com/findy-network/findy-agent-vault/graph/model"
)

func (r *queryResolver) User(_ context.Context) (*model.User, error) {
	return data.State.User.ToNode(), nil
}
