package resolver

import (
	"context"

	"github.com/findy-network/findy-agent-vault/db/store"
	"github.com/findy-network/findy-agent-vault/graph/model"
	"github.com/findy-network/findy-agent-vault/utils"
	"github.com/lainio/err2"
)

func (r *jobResolver) output(ctx context.Context, obj *model.Job) (o *model.JobOutput, err error) {
	defer err2.Return(&err)

	// TODO: store agent data to context?
	agent, err := store.GetAgent(ctx, r.db)
	err2.Check(err)

	utils.LogMed().Infof(
		"jobResolver:Output for tenant %s, event: %s",
		agent.ID,
		obj.ID,
	)

	output, err := r.db.GetJobOutput(obj.ID, agent.ID, obj.Protocol)
	err2.Check(err)

	return output.ToEdges(), nil
}
