package eventconn

import (
	"context"

	"github.com/findy-network/findy-agent-vault/db/store"
	"github.com/findy-network/findy-agent-vault/graph/model"
	"github.com/findy-network/findy-agent-vault/resolver/agent"
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

func (r *Resolver) TotalCount(ctx context.Context, obj *model.EventConnection) (c int, err error) {
	defer err2.Return(&err)

	tenant, err := r.GetAgent(ctx)
	err2.Check(err)

	utils.LogLow().Infof(
		"eventConnectionResolver:TotalCount for tenant %s, connection: %v",
		tenant.ID,
		obj.ConnectionID,
	)
	count, err := r.db.GetEventCount(tenant.ID, obj.ConnectionID)
	err2.Check(err)

	return count, nil
}
