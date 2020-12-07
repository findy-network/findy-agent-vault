package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/findy-network/findy-agent-vault/graph/generated"
	"github.com/findy-network/findy-agent-vault/graph/model"
	"github.com/findy-network/findy-agent-vault/utils"
)

func (r *basicMessageResolver) Connection(ctx context.Context, obj *model.BasicMessage) (*model.Pairwise, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *credentialResolver) Connection(ctx context.Context, obj *model.Credential) (*model.Pairwise, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *eventResolver) Job(ctx context.Context, obj *model.Event) (*model.JobEdge, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *eventResolver) Connection(ctx context.Context, obj *model.Event) (*model.Pairwise, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *jobResolver) Output(ctx context.Context, obj *model.Job) (*model.JobOutput, error) {
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

func (r *mutationResolver) Resume(ctx context.Context, input model.ResumeJobInput) (*model.Response, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AddRandomEvent(ctx context.Context) (bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AddRandomMessage(ctx context.Context) (bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AddRandomCredential(ctx context.Context) (bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AddRandomProof(ctx context.Context) (bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *pairwiseResolver) Messages(ctx context.Context, obj *model.Pairwise, after *string, before *string, first *int, last *int) (*model.BasicMessageConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *pairwiseResolver) Credentials(ctx context.Context, obj *model.Pairwise, after *string, before *string, first *int, last *int) (*model.CredentialConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *pairwiseResolver) Proofs(ctx context.Context, obj *model.Pairwise, after *string, before *string, first *int, last *int) (*model.ProofConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *pairwiseResolver) Jobs(ctx context.Context, obj *model.Pairwise, after *string, before *string, first *int, last *int, completed *bool) (*model.JobConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *pairwiseResolver) Events(ctx context.Context, obj *model.Pairwise, after *string, before *string, first *int, last *int) (*model.EventConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *proofResolver) Connection(ctx context.Context, obj *model.Proof) (*model.Pairwise, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Connections(ctx context.Context, after *string, before *string, first *int, last *int) (*model.PairwiseConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Connection(ctx context.Context, id string) (*model.Pairwise, error) {
	utils.LogMed().Info("queryResolver:Connection, id: ", id)

	agentID, _ := utils.ParseUser(ctx)

	conn, err := r.db.GetConnection(id, agentID)
	if err != nil {
		return nil, err
	}
	return conn.ToNode(), nil
}

func (r *queryResolver) Message(ctx context.Context, id string) (*model.BasicMessage, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Credential(ctx context.Context, id string) (*model.Credential, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Credentials(ctx context.Context, after *string, before *string, first *int, last *int) (*model.CredentialConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Proof(ctx context.Context, id string) (*model.Proof, error) {
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
