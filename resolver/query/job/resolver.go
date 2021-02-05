package job

import (
	"context"

	"github.com/findy-network/findy-agent-vault/db/store"
	graph "github.com/findy-network/findy-agent-vault/graph/model"
	"github.com/findy-network/findy-agent-vault/resolver/query/agent"
	"github.com/findy-network/findy-agent-vault/utils"
	"github.com/lainio/err2"
)

type Resolver struct {
	db store.DB
	*agent.Resolver
}

func NewResolver(db store.DB, agentResolver *agent.Resolver) *Resolver {
	return &Resolver{db, agentResolver}
}

func (r *Resolver) Output(ctx context.Context, obj *graph.Job) (o *graph.JobOutput, err error) {
	defer err2.Return(&err)

	tenant, err := r.GetAgent(ctx)
	err2.Check(err)

	utils.LogLow().Infof(
		"jobResolver:Output for tenant %s, event: %s",
		tenant.ID,
		obj.ID,
	)

	output, err := r.db.GetJobOutput(obj.ID, tenant.ID, obj.Protocol)
	err2.Check(err)

	return output.ToEdges(), nil
}
