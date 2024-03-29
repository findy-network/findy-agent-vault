package proof

import (
	"context"

	"github.com/findy-network/findy-agent-vault/db/store"
	"github.com/findy-network/findy-agent-vault/graph/model"
	"github.com/findy-network/findy-agent-vault/resolver/query/agent"
	"github.com/findy-network/findy-agent-vault/utils"
	"github.com/golang/glog"
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

func (r *Resolver) Connection(ctx context.Context, obj *model.Proof) (c *model.Pairwise, err error) {
	defer err2.Handle(&err)

	tenant := try.To1(r.GetAgent(ctx))

	utils.LogLow().Infof(
		"proofResolver:Connection for tenant %s, proof: %s",
		tenant.ID,
		obj.ID,
	)

	connection := try.To1(r.db.GetConnectionForProof(obj.ID, tenant.ID))

	return connection.ToNode(), nil
}

func (r *Resolver) Provable(ctx context.Context, obj *model.Proof) (res *model.Provable, err error) {
	defer err2.Handle(&err)

	tenant := try.To1(r.GetAgent(ctx))

	utils.LogLow().Infof(
		"proofResolver:Provable for tenant %s, proof : %s",
		tenant.ID,
		obj.ID,
	)

	res = &model.Provable{ID: obj.ID}

	provable := false
	// provable only if not accepted yet
	if obj.Role == model.ProofRoleProver && obj.ApprovedMs == nil && obj.VerifiedMs == nil {
		res.Attributes, err = r.db.SearchCredentials(tenant.ID, obj.Attributes)
		if err == nil {
			provable = true
			for _, attr := range res.Attributes {
				if len(attr.Credentials) == 0 {
					provable = false
					break
				}
			}
		} else {
			glog.Warningf("Encountered error when searching credentials: %s %s", tenant.ID, err)
		}
	}

	res.Provable = provable
	return res, nil
}
