package message

import (
	"context"

	"github.com/findy-network/findy-agent-vault/db/store"
	"github.com/findy-network/findy-agent-vault/graph/model"
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

func (r *Resolver) Connection(ctx context.Context, obj *model.BasicMessage) (c *model.Pairwise, err error) {
	defer err2.Handle(&err)

	tenant := try.To1(r.GetAgent(ctx))

	utils.LogLow().Infof(
		"basicMessageResolver:Connection for tenant %s, message: %s",
		tenant.ID,
		obj.ID,
	)

	connection := try.To1(r.db.GetConnectionForMessage(obj.ID, tenant.ID))

	return connection.ToNode(), nil
}
