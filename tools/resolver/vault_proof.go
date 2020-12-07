package resolver

import (
	"context"
	"fmt"

	"github.com/findy-network/findy-agent-vault/graph/model"
	data "github.com/findy-network/findy-agent-vault/tools/data/model"
	"github.com/findy-network/findy-agent-vault/tools/faker"
	"github.com/findy-network/findy-agent-vault/tools/tools"
	"github.com/findy-network/findy-agent-vault/utils"
	"github.com/lainio/err2"
)

func (r *proofResolver) Connection(ctx context.Context, obj *model.Proof) (p *model.Pairwise, err error) {
	utils.LogMed().Info("proofResolver:Connection, id: ", obj.ID)
	defer err2.Return(&err)

	if connectionID := state.Proofs().ProofPairwiseID(obj.ID); connectionID != nil {
		return r.Query().Connection(ctx, *connectionID)
	}

	err = fmt.Errorf("pairwise for proof id %s was not found", obj.ID)
	return
}

func (r *queryResolver) Proof(ctx context.Context, id string) (node *model.Proof, err error) {
	utils.LogMed().Info("queryResolver:Proof, id: ", id)

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

	utils.LogLow().Infof("Proofs: returning proofs between %d and %d", afterIndex, beforeIndex)

	return items.ProofConnection(afterIndex, beforeIndex), nil
}

func (r *mutationResolver) AddRandomProof(ctx context.Context) (ok bool, err error) {
	utils.LogMed().Info("mutationResolver:AddRandomProof ")
	defer err2.Return(&err)

	proofs, err := faker.FakeProofs(1)
	err2.Check(err)

	proof := proofs[0]
	r.listener.AddProof(
		proof.PairwiseID,
		proof.ID,
		proof.Role,
		proof.Attributes,
		proof.InitiatedByUs,
	)
	currentTime := tools.CurrentTimeMs()
	r.listener.UpdateProof(proof.PairwiseID, proof.ID, &currentTime, &currentTime, nil)

	ok = true

	return
}
