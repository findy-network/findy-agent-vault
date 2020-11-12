package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/findy-network/findy-agent-vault/graph/generated"
	"github.com/findy-network/findy-agent-vault/graph/model"
)

func (r *basicMessageResolver) Connection(ctx context.Context, obj *model.BasicMessage) (*model.Pairwise, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) MarkEventRead(ctx context.Context, input model.MarkReadInput) (*model.Event, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) Invite(ctx context.Context) (*model.InvitationResponse, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) Connect(ctx context.Context, input model.ConnectInput) (*model.Response, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) SendMessage(ctx context.Context, input model.MessageInput) (*model.Response, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AcceptOffer(ctx context.Context, input model.Offer) (*model.Response, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AcceptRequest(ctx context.Context, input model.Request) (*model.Response, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AddRandomEvent(ctx context.Context) (bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *pairwiseResolver) Messages(ctx context.Context, obj *model.Pairwise, after *string, before *string, first *int, last *int) (*model.BasicMessageConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Connections(ctx context.Context, after *string, before *string, first *int, last *int) (*model.PairwiseConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Connection(ctx context.Context, id string) (*model.Pairwise, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Message(ctx context.Context, id string) (*model.BasicMessage, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Events(ctx context.Context, after *string, before *string, first *int, last *int) (*model.EventConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Event(ctx context.Context, id string) (*model.Event, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Jobs(ctx context.Context, after *string, before *string, first *int, last *int, completed *bool) (*model.JobConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Job(ctx context.Context, id string) (*model.Job, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) User(ctx context.Context) (*model.User, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *subscriptionResolver) EventAdded(ctx context.Context) (<-chan *model.EventEdge, error) {
	panic(fmt.Errorf("not implemented"))
}

// BasicMessage returns generated.BasicMessageResolver implementation.
func (r *Resolver) BasicMessage() generated.BasicMessageResolver { return &basicMessageResolver{r} }

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Pairwise returns generated.PairwiseResolver implementation.
func (r *Resolver) Pairwise() generated.PairwiseResolver { return &pairwiseResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// Subscription returns generated.SubscriptionResolver implementation.
func (r *Resolver) Subscription() generated.SubscriptionResolver { return &subscriptionResolver{r} }

type basicMessageResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
type pairwiseResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }
