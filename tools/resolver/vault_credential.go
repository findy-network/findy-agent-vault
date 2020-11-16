package resolver

import (
	"context"
	"fmt"

	"github.com/findy-network/findy-agent-vault/graph/model"
)

func (r *credentialResolver) Connection(ctx context.Context, obj *model.Credential) (*model.Pairwise, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Credential(ctx context.Context, id string) (*model.Credential, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Credentials(ctx context.Context, after *string, before *string, first *int, last *int) (*model.CredentialConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *pairwiseResolver) Credentials(ctx context.Context, obj *model.Pairwise, after *string, before *string, first *int, last *int) (*model.CredentialConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

