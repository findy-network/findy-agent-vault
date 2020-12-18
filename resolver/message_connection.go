package resolver

import (
	"context"

	"github.com/findy-network/findy-agent-vault/db/store"
	"github.com/findy-network/findy-agent-vault/graph/model"
	"github.com/findy-network/findy-agent-vault/utils"
	"github.com/lainio/err2"
)

func (r *basicMessageConnectionResolver) totalCount(ctx context.Context, obj *model.BasicMessageConnection) (c int, err error) {
	defer err2.Return(&err)

	// TODO: store agent data to context?
	agent, err := store.GetAgent(ctx, r.db)
	err2.Check(err)

	utils.LogMed().Infof(
		"MessageConnectionResolver:TotalCount for tenant %s, connection: %v",
		agent.ID,
		obj.ConnectionID,
	)
	count, err := r.db.GetMessageCount(agent.ID, obj.ConnectionID)
	err2.Check(err)

	return count, nil
}
