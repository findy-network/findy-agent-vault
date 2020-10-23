package resolver

import (
	"context"
	"fmt"

	"github.com/findy-network/findy-agent-api/tools/faker"

	"github.com/findy-network/findy-agent-api/tools/data"

	"github.com/findy-network/findy-agent-api/graph/generated"
	"github.com/findy-network/findy-agent-api/graph/model"
)

const (
	logLevelMedium = 2
)

func InitResolver() {
	data.InitState()
	faker.InitFaker()
	initEvents()
}

type Resolver struct{}

func (r *mutationResolver) Invite(_ context.Context) (*model.Response, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) Connect(_ context.Context, _ model.Invitation) (*model.Response, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) SendMessage(_ context.Context) (*model.Response, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AcceptOffer(_ context.Context, _ model.Offer) (*model.Response, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AcceptRequest(_ context.Context, _ model.Request) (*model.Response, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Event(_ context.Context, id string) (*model.Event, error) {
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
