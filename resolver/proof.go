package resolver

import (
	"context"

	"github.com/findy-network/findy-agent-vault/graph/model"
	"github.com/findy-network/findy-agent-vault/utils"
	"github.com/lainio/err2"
)

func (r *proofResolver) connection(ctx context.Context, obj *model.Proof) (c *model.Pairwise, err error) {
	defer err2.Return(&err)

	// TODO: store agent data to context?
	agent, err := r.getAgent(ctx)
	err2.Check(err)

	utils.LogMed().Infof(
		"proofResolver:Connection for tenant %s, proof: %s",
		agent.ID,
		obj.ID,
	)

	connection, err := r.db.GetConnectionForProof(obj.ID, agent.ID)
	err2.Check(err)

	return connection.ToNode(), nil
}
