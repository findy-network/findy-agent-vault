package resolver

import (
	"context"

	"github.com/findy-network/findy-agent-vault/graph/model"
	"github.com/findy-network/findy-agent-vault/utils"
	"github.com/lainio/err2"
)

func (r *proofConnectionResolver) totalCount(ctx context.Context, obj *model.ProofConnection) (c int, err error) {
	defer err2.Return(&err)

	// TODO: store agent data to context?
	agent, err := r.getAgent(ctx)
	err2.Check(err)

	utils.LogMed().Infof(
		"proofConnectionResolver:TotalCount for tenant %s, connection: %v",
		agent.ID,
		obj.ConnectionID,
	)

	count, err := r.db.GetProofCount(agent.ID, obj.ConnectionID)
	err2.Check(err)

	return count, nil
}
