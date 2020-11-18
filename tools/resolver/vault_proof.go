package resolver

import (
	"context"
	"fmt"

	"github.com/findy-network/findy-agent-vault/graph/model"
	data "github.com/findy-network/findy-agent-vault/tools/data/model"
	"github.com/golang/glog"
	"github.com/lainio/err2"
)

func (r *proofResolver) Connection(ctx context.Context, obj *model.Proof) (p *model.Pairwise, err error) {
	glog.V(logLevelMedium).Info("proofResolver:Connection, id: ", obj.ID)
	defer err2.Return(&err)

	if connectionID := state.Proofs().ProofPairwiseID(obj.ID); connectionID != nil {
		return r.Query().Connection(ctx, *connectionID)
	}

	err = fmt.Errorf("pairwise for proof id %s was not found", obj.ID)
	return
}

func (r *queryResolver) Proof(ctx context.Context, id string) (node *model.Proof, err error) {
	glog.V(logLevelMedium).Info("queryResolver:Proof, id: ", id)

	items := state.Proofs()
	edge := items.ProofForID(id)
	if edge == nil {
		err = fmt.Errorf("connection for id %s was not found", id)
	} else {
		node = edge.Node
	}
	return
}

func (r *pairwiseResolver) Proofs(
	ctx context.Context,
	obj *model.Pairwise,
	after, before *string,
	first, last *int,
) (c *model.ProofConnection, err error) {
	defer err2.Return(&err)
	pagination := &PaginationParams{
		first:  first,
		last:   last,
		after:  after,
		before: before,
	}
	logPaginationRequest("pairwiseResolver:proofs", pagination)

	items := state.Proofs()
	items = &data.ProofItems{Items: items.Filter(func(item data.APIObject) data.APIObject {
		c := item.Proof()
		if c.VerifiedMs != nil && c.PairwiseID == obj.ID {
			return c.Copy()
		}
		return nil
	})}

	afterIndex, beforeIndex, err := pick(items.Objects(), pagination)
	err2.Check(err)

	glog.V(logLevelLow).Infof("Proofs: returning proofs between %d and %d", afterIndex, beforeIndex)

	return items.ProofConnection(afterIndex, beforeIndex), nil
}
