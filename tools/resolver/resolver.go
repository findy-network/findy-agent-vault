package resolver

import (
	"context"
	"fmt"

	"github.com/findy-network/findy-agent-vault/tools/data"

	"github.com/findy-network/findy-agent-vault/agency"
	"github.com/findy-network/findy-agent-vault/graph/generated"
	"github.com/findy-network/findy-agent-vault/graph/model"
)

const (
	logLevelMedium = 2
	logLevelLow    = 3
)

type agencyListener struct{}

var state *data.Data

func InitResolver() *Resolver {
	agency.Instance.Init(&agencyListener{})
	state = data.InitState()
	initEvents()
	return &Resolver{}
}

type Resolver struct{}

func (r *mutationResolver) SendMessage(_ context.Context) (*model.Response, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AcceptOffer(_ context.Context, _ model.Offer) (*model.Response, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AcceptRequest(_ context.Context, _ model.Request) (*model.Response, error) {
	panic(fmt.Errorf("not implemented"))
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// Subscription returns generated.SubscriptionResolver implementation.
func (r *Resolver) Subscription() generated.SubscriptionResolver { return &subscriptionResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }
