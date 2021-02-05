package event

import (
	"context"

	"github.com/findy-network/findy-agent-vault/db/store"
	"github.com/findy-network/findy-agent-vault/graph/model"
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

func (r *Resolver) Connection(ctx context.Context, obj *model.Event) (c *model.Pairwise, err error) {
	defer err2.Return(&err)

	tenant, err := r.GetAgent(ctx)
	err2.Check(err)

	utils.LogLow().Infof(
		"eventResolver:Connection for tenant %s, event: %s",
		tenant.ID,
		obj.ID,
	)

	connection, err := r.db.GetConnectionForEvent(obj.ID, tenant.ID)
	err2.Check(err)

	return connection.ToNode(), nil
}

func (r *Resolver) Job(ctx context.Context, obj *model.Event) (j *model.JobEdge, err error) {
	defer err2.Return(&err)

	tenant, err := r.GetAgent(ctx)
	err2.Check(err)

	utils.LogLow().Infof(
		"eventResolver:Job for tenant %s, event: %s",
		tenant.ID,
		obj.ID,
	)

	job, err := r.db.GetJobForEvent(obj.ID, tenant.ID)
	err2.Check(err)

	return job.ToEdge(), nil
}
