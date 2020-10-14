package resolver

import (
	"context"
	"fmt"

	"github.com/findy-network/findy-agent-api/graph/generated"
	"github.com/findy-network/findy-agent-api/graph/model"
)

func init() {
	initEvents()
}

type Resolver struct{}

func (r *mutationResolver) Invite(ctx context.Context) (*model.Response, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) Connect(ctx context.Context, input model.Invitation) (*model.Response, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) SendMessage(ctx context.Context) (*model.Response, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AcceptOffer(ctx context.Context, input model.Offer) (*model.Response, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AcceptRequest(ctx context.Context, input model.Request) (*model.Response, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Event(ctx context.Context, id string) (*model.Event, error) {
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
