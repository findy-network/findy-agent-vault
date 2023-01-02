package job

import (
	"context"

	"github.com/findy-network/findy-agent-vault/db/store"
	graph "github.com/findy-network/findy-agent-vault/graph/model"
	"github.com/findy-network/findy-agent-vault/resolver/query/agent"
	"github.com/findy-network/findy-agent-vault/utils"
	"github.com/lainio/err2"
	"github.com/lainio/err2/try"
)

type Resolver struct {
	db store.DB
	*agent.Resolver
}

func NewResolver(db store.DB, agentResolver *agent.Resolver) *Resolver {
	return &Resolver{db, agentResolver}
}

func (r *Resolver) Output(ctx context.Context, obj *graph.Job) (o *graph.JobOutput, err error) {
	defer err2.Handle(&err)

	tenant := try.To1(r.GetAgent(ctx))

	utils.LogLow().Infof(
		"jobResolver:Output for tenant %s, event: %s",
		tenant.ID,
		obj.ID,
	)

	output := try.To1(r.db.GetJobOutput(obj.ID, tenant.ID, obj.Protocol))

	return output.ToEdges(), nil
}
