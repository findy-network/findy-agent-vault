package resolver

import (
	"github.com/findy-network/findy-agent-vault/tools/data"

	"github.com/findy-network/findy-agent-vault/agency"
	"github.com/findy-network/findy-agent-vault/graph/generated"
)

type agencyListener struct{}

var state *data.Data

func InitResolver(skipFake bool) *Resolver {
	listener := &agencyListener{}
	agency.Instance.Init(listener)
	state = data.InitState(skipFake)
	initEvents()
	return &Resolver{listener}
}

type Resolver struct {
	listener *agencyListener
}

// BasicMessage returns generated.BasicMessageResolver implementation.
func (r *Resolver) BasicMessage() generated.BasicMessageResolver { return &basicMessageResolver{r} }

// Credential returns generated.CredentialResolver implementation.
func (r *Resolver) Credential() generated.CredentialResolver { return &credentialResolver{r} }

// Event returns generated.EventResolver implementation.
func (r *Resolver) Event() generated.EventResolver { return &eventResolver{r} }

// Job returns generated.JobResolver implementation.
func (r *Resolver) Job() generated.JobResolver { return &jobResolver{r} }

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Pairwise returns generated.PairwiseResolver implementation.
func (r *Resolver) Pairwise() generated.PairwiseResolver { return &pairwiseResolver{r} }

// Proof returns generated.ProofResolver implementation.
func (r *Resolver) Proof() generated.ProofResolver { return &proofResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// Subscription returns generated.SubscriptionResolver implementation.
func (r *Resolver) Subscription() generated.SubscriptionResolver { return &subscriptionResolver{r} }

type basicMessageResolver struct{ *Resolver }
type credentialResolver struct{ *Resolver }
type eventResolver struct{ *Resolver }
type jobResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
type pairwiseResolver struct{ *Resolver }
type proofResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }
