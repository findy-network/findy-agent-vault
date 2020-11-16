package resolver

import (
	"context"
	"fmt"

	"github.com/findy-network/findy-agent-vault/graph/model"
)

func (r *queryResolver) Proof(ctx context.Context, id string) (*model.Proof, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *proofResolver) Connection(ctx context.Context, obj *model.Proof) (*model.Pairwise, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *pairwiseResolver) Proofs(ctx context.Context, obj *model.Pairwise, after *string, before *string, first *int, last *int) (*model.ProofConnection, error) {
	panic(fmt.Errorf("not implemented"))
}
