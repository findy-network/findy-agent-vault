package resolver

import (
	"context"

	"github.com/findy-network/findy-agent-vault/graph/model"
	"github.com/findy-network/findy-agent-vault/utils"
	"github.com/lainio/err2"
)

func (r *pairwiseConnectionResolver) totalCount(ctx context.Context, obj *model.PairwiseConnection) (c int, err error) {
	defer err2.Return(&err)

	// TODO: store agent data to context?
	agent, err := r.getAgent(ctx)
	err2.Check(err)

	utils.LogLow().Infof("pairwiseConnectionResolver:TotalCount for tenant %s", agent.ID)

	count, err := r.db.GetConnectionCount(agent.ID)
	err2.Check(err)

	return count, nil
}
