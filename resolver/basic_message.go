package resolver

import (
	"context"

	"github.com/findy-network/findy-agent-vault/db/store"
	"github.com/findy-network/findy-agent-vault/graph/model"
	"github.com/findy-network/findy-agent-vault/utils"
	"github.com/lainio/err2"
)

func (r *basicMessageResolver) connection(ctx context.Context, obj *model.BasicMessage) (c *model.Pairwise, err error) {
	defer err2.Return(&err)

	// TODO: store agent data to context?
	agent, err := store.GetAgent(ctx, r.db)
	err2.Check(err)

	utils.LogMed().Infof(
		"basicMessageResolver:Connection for tenant %s, message: %s",
		agent.ID,
		obj.ID,
	)

	connection, err := r.db.GetConnectionForMessage(obj.ID, agent.ID)
	err2.Check(err)

	return connection.ToNode(), nil
}
