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
	return r.resolvers.message.Connection(ctx, obj)
}

func (r *basicMessageConnectionResolver) TotalCount(ctx context.Context, obj *model.BasicMessageConnection) (int, error) {
	return r.resolvers.messageConnection.TotalCount(ctx, obj)
}

func (r *credentialResolver) Connection(ctx context.Context, obj *model.Credential) (*model.Pairwise, error) {
	return r.resolvers.credential.Connection(ctx, obj)
}

func (r *credentialConnectionResolver) TotalCount(ctx context.Context, obj *model.CredentialConnection) (int, error) {
	return r.resolvers.credentialConnection.TotalCount(ctx, obj)
}

func (r *eventResolver) Job(ctx context.Context, obj *model.Event) (*model.JobEdge, error) {
	return r.resolvers.event.Job(ctx, obj)
}

func (r *eventResolver) Connection(ctx context.Context, obj *model.Event) (*model.Pairwise, error) {
	return r.resolvers.event.Connection(ctx, obj)
}

func (r *eventConnectionResolver) TotalCount(ctx context.Context, obj *model.EventConnection) (int, error) {
	return r.resolvers.eventConnection.TotalCount(ctx, obj)
}

func (r *jobResolver) Output(ctx context.Context, obj *model.Job) (*model.JobOutput, error) {
	return r.resolvers.job.Output(ctx, obj)
}

func (r *jobConnectionResolver) TotalCount(ctx context.Context, obj *model.JobConnection) (int, error) {
	return r.resolvers.jobConnection.TotalCount(ctx, obj)
}

func (r *mutationResolver) MarkEventRead(ctx context.Context, input model.MarkReadInput) (*model.Event, error) {
	return r.resolvers.mutation.MarkEventRead(ctx, input)
}

func (r *mutationResolver) Invite(ctx context.Context) (*model.InvitationResponse, error) {
	return r.resolvers.mutation.Invite(ctx)
}

func (r *mutationResolver) Connect(ctx context.Context, input model.ConnectInput) (*model.Response, error) {
	return r.resolvers.mutation.Connect(ctx, input)
}

func (r *mutationResolver) SendMessage(ctx context.Context, input model.MessageInput) (*model.Response, error) {
	return r.resolvers.mutation.SendMessage(ctx, input)
}

func (r *mutationResolver) Resume(ctx context.Context, input model.ResumeJobInput) (*model.Response, error) {
	return r.resolvers.mutation.Resume(ctx, input)
}

func (r *mutationResolver) AddRandomEvent(ctx context.Context) (bool, error) {
	return r.resolvers.playground.AddRandomEvent(ctx)
}

func (r *mutationResolver) AddRandomMessage(ctx context.Context) (bool, error) {
	return r.resolvers.playground.AddRandomMessage(ctx)
}

func (r *mutationResolver) AddRandomCredential(ctx context.Context) (bool, error) {
	return r.resolvers.playground.AddRandomCredential(ctx)
}

func (r *mutationResolver) AddRandomProof(ctx context.Context) (bool, error) {
	return r.resolvers.playground.AddRandomProof(ctx)
}

func (r *pairwiseResolver) Messages(ctx context.Context, obj *model.Pairwise, after *string, before *string, first *int, last *int) (*model.BasicMessageConnection, error) {
	return r.resolvers.pairwise.Messages(ctx, obj, after, before, first, last)
}

func (r *pairwiseResolver) Credentials(ctx context.Context, obj *model.Pairwise, after *string, before *string, first *int, last *int) (*model.CredentialConnection, error) {
	return r.resolvers.pairwise.Credentials(ctx, obj, after, before, first, last)
}

func (r *pairwiseResolver) Proofs(ctx context.Context, obj *model.Pairwise, after *string, before *string, first *int, last *int) (*model.ProofConnection, error) {
	return r.resolvers.pairwise.Proofs(ctx, obj, after, before, first, last)
}

func (r *pairwiseResolver) Jobs(ctx context.Context, obj *model.Pairwise, after *string, before *string, first *int, last *int, completed *bool) (*model.JobConnection, error) {
	return r.resolvers.pairwise.Jobs(ctx, obj, after, before, first, last, completed)
}

func (r *pairwiseResolver) Events(ctx context.Context, obj *model.Pairwise, after *string, before *string, first *int, last *int) (*model.EventConnection, error) {
	return r.resolvers.pairwise.Events(ctx, obj, after, before, first, last)
}

func (r *pairwiseConnectionResolver) TotalCount(ctx context.Context, obj *model.PairwiseConnection) (int, error) {
	return r.resolvers.pairwiseConnection.TotalCount(ctx, obj)
}

func (r *proofResolver) Connection(ctx context.Context, obj *model.Proof) (*model.Pairwise, error) {
	return r.resolvers.proof.Connection(ctx, obj)
}

func (r *proofAttributeResolver) Credentials(ctx context.Context, obj *model.ProofAttribute) ([]*model.CredentialMatch, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *proofConnectionResolver) TotalCount(ctx context.Context, obj *model.ProofConnection) (int, error) {
	return r.resolvers.proofConnection.TotalCount(ctx, obj)
}

func (r *queryResolver) Connections(ctx context.Context, after *string, before *string, first *int, last *int) (*model.PairwiseConnection, error) {
	return r.resolvers.query.Connections(ctx, after, before, first, last)
}

func (r *queryResolver) Connection(ctx context.Context, id string) (*model.Pairwise, error) {
	return r.resolvers.query.Connection(ctx, id)
}

func (r *queryResolver) Message(ctx context.Context, id string) (*model.BasicMessage, error) {
	return r.resolvers.query.Message(ctx, id)
}

func (r *queryResolver) Credential(ctx context.Context, id string) (*model.Credential, error) {
	return r.resolvers.query.Credential(ctx, id)
}

func (r *queryResolver) Credentials(ctx context.Context, after *string, before *string, first *int, last *int) (*model.CredentialConnection, error) {
	return r.resolvers.query.Credentials(ctx, after, before, first, last)
}

func (r *queryResolver) Proof(ctx context.Context, id string) (*model.Proof, error) {
	return r.resolvers.query.Proof(ctx, id)
}

func (r *queryResolver) Events(ctx context.Context, after *string, before *string, first *int, last *int) (*model.EventConnection, error) {
	return r.resolvers.query.Events(ctx, after, before, first, last)
}

func (r *queryResolver) Event(ctx context.Context, id string) (*model.Event, error) {
	return r.resolvers.query.Event(ctx, id)
}

func (r *queryResolver) Jobs(ctx context.Context, after *string, before *string, first *int, last *int, completed *bool) (*model.JobConnection, error) {
	return r.resolvers.query.Jobs(ctx, after, before, first, last, completed)
}

func (r *queryResolver) Job(ctx context.Context, id string) (*model.Job, error) {
	return r.resolvers.query.Job(ctx, id)
}

func (r *queryResolver) User(ctx context.Context) (*model.User, error) {
	return r.resolvers.query.User(ctx)
}

func (r *queryResolver) Endpoint(ctx context.Context, payload string) (*model.InvitationResponse, error) {
	return r.resolvers.query.Endpoint(ctx, payload)
}

func (r *subscriptionResolver) EventAdded(ctx context.Context) (<-chan *model.EventEdge, error) {
	return r.updater.EventAdded(ctx)
}

// BasicMessage returns generated.BasicMessageResolver implementation.
func (r *Resolver) BasicMessage() generated.BasicMessageResolver { return &basicMessageResolver{r} }

// BasicMessageConnection returns generated.BasicMessageConnectionResolver implementation.
func (r *Resolver) BasicMessageConnection() generated.BasicMessageConnectionResolver {
	return &basicMessageConnectionResolver{r}
}

// Credential returns generated.CredentialResolver implementation.
func (r *Resolver) Credential() generated.CredentialResolver { return &credentialResolver{r} }

// CredentialConnection returns generated.CredentialConnectionResolver implementation.
func (r *Resolver) CredentialConnection() generated.CredentialConnectionResolver {
	return &credentialConnectionResolver{r}
}

// Event returns generated.EventResolver implementation.
func (r *Resolver) Event() generated.EventResolver { return &eventResolver{r} }

// EventConnection returns generated.EventConnectionResolver implementation.
func (r *Resolver) EventConnection() generated.EventConnectionResolver {
	return &eventConnectionResolver{r}
}

// Job returns generated.JobResolver implementation.
func (r *Resolver) Job() generated.JobResolver { return &jobResolver{r} }

// JobConnection returns generated.JobConnectionResolver implementation.
func (r *Resolver) JobConnection() generated.JobConnectionResolver { return &jobConnectionResolver{r} }

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Pairwise returns generated.PairwiseResolver implementation.
func (r *Resolver) Pairwise() generated.PairwiseResolver { return &pairwiseResolver{r} }

// PairwiseConnection returns generated.PairwiseConnectionResolver implementation.
func (r *Resolver) PairwiseConnection() generated.PairwiseConnectionResolver {
	return &pairwiseConnectionResolver{r}
}

// Proof returns generated.ProofResolver implementation.
func (r *Resolver) Proof() generated.ProofResolver { return &proofResolver{r} }

// ProofAttribute returns generated.ProofAttributeResolver implementation.
func (r *Resolver) ProofAttribute() generated.ProofAttributeResolver {
	return &proofAttributeResolver{r}
}

// ProofConnection returns generated.ProofConnectionResolver implementation.
func (r *Resolver) ProofConnection() generated.ProofConnectionResolver {
	return &proofConnectionResolver{r}
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// Subscription returns generated.SubscriptionResolver implementation.
func (r *Resolver) Subscription() generated.SubscriptionResolver { return &subscriptionResolver{r} }

type basicMessageResolver struct{ *Resolver }
type basicMessageConnectionResolver struct{ *Resolver }
type credentialResolver struct{ *Resolver }
type credentialConnectionResolver struct{ *Resolver }
type eventResolver struct{ *Resolver }
type eventConnectionResolver struct{ *Resolver }
type jobResolver struct{ *Resolver }
type jobConnectionResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
type pairwiseResolver struct{ *Resolver }
type pairwiseConnectionResolver struct{ *Resolver }
type proofResolver struct{ *Resolver }
type proofAttributeResolver struct{ *Resolver }
type proofConnectionResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }
