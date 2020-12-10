package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/findy-network/findy-agent-vault/db/db"
	"github.com/findy-network/findy-agent-vault/graph/generated"
	"github.com/findy-network/findy-agent-vault/graph/model"
	"github.com/findy-network/findy-agent-vault/paginator"
	"github.com/findy-network/findy-agent-vault/utils"
	"github.com/lainio/err2"
)

func (r *basicMessageResolver) Connection(ctx context.Context, obj *model.BasicMessage) (*model.Pairwise, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *credentialResolver) Connection(ctx context.Context, obj *model.Credential) (*model.Pairwise, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *credentialConnectionResolver) TotalCount(ctx context.Context, obj *model.CredentialConnection) (int, error) {
	var err error
	defer err2.Return(&err)

	// TODO: store agent data to context?
	agent, err := db.GetAgent(ctx, r.db)
	err2.Check(err)

	utils.LogMed().Infof("credentialConnectionResolver:TotalCount for tenant %s", agent.ID)

	count, err := r.db.GetCredentialCount(agent.ID)
	err2.Check(err)

	return count, nil
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
	var err error
	defer err2.Return(&err)

	agent, err := db.GetAgent(ctx, r.db)
	err2.Check(err)

	utils.LogMed().Infof("pairwiseResolver:Credentials for tenant: %s, connection %s", agent.ID, obj.ID)

	batch, err := paginator.Validate("pairwiseResolver:Credentials", &paginator.Params{
		First:  first,
		Last:   last,
		After:  after,
		Before: before,
	})
	err2.Check(err)

	res, err := r.db.GetConnectionCredentials(batch, agent.ID, obj.ID)
	err2.Check(err)

	return res.ToConnection(&obj.ID), nil
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

func (r *pairwiseConnectionResolver) TotalCount(ctx context.Context, obj *model.PairwiseConnection) (int, error) {
	var err error
	defer err2.Return(&err)

	// TODO: store agent data to context?
	agent, err := db.GetAgent(ctx, r.db)
	err2.Check(err)

	utils.LogMed().Infof("pairwiseConnectionResolver:TotalCount for tenant %s", agent.ID)

	count, err := r.db.GetConnectionCount(agent.ID)
	err2.Check(err)

	return count, nil
}

func (r *proofResolver) Connection(ctx context.Context, obj *model.Proof) (*model.Pairwise, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Connections(ctx context.Context, after *string, before *string, first *int, last *int) (*model.PairwiseConnection, error) {
	var err error
	defer err2.Return(&err)

	agent, err := db.GetAgent(ctx, r.db)
	err2.Check(err)

	utils.LogMed().Info("queryResolver:Connections for tenant: ", agent.ID)

	batch, err := paginator.Validate("queryResolver:Connections", &paginator.Params{
		First:  first,
		Last:   last,
		After:  after,
		Before: before,
	})
	err2.Check(err)

	res, err := r.db.GetConnections(batch, agent.ID)
	err2.Check(err)

	return res.ToConnection(), nil
}

func (r *queryResolver) Connection(ctx context.Context, id string) (*model.Pairwise, error) {
	var err error
	defer err2.Return(&err)

	agent, err := db.GetAgent(ctx, r.db)
	err2.Check(err)

	utils.LogMed().Infof("queryResolver:Connection id: %s for tenant %s", id, agent.ID)

	conn, err := r.db.GetConnection(id, agent.ID)
	err2.Check(err)

	return conn.ToNode(), nil
}

func (r *queryResolver) Message(ctx context.Context, id string) (*model.BasicMessage, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Credential(ctx context.Context, id string) (*model.Credential, error) {
	var err error
	defer err2.Return(&err)

	agent, err := db.GetAgent(ctx, r.db)
	err2.Check(err)

	utils.LogMed().Infof("queryResolver:Credential id: %s for tenant %s", id, agent.ID)

	cred, err := r.db.GetCredential(id, agent.ID)
	err2.Check(err)

	return cred.ToNode(), nil
}

func (r *queryResolver) Credentials(ctx context.Context, after *string, before *string, first *int, last *int) (*model.CredentialConnection, error) {
	var err error
	defer err2.Return(&err)

	agent, err := db.GetAgent(ctx, r.db)
	err2.Check(err)

	utils.LogMed().Info("queryResolver:Credentials for tenant: ", agent.ID)

	batch, err := paginator.Validate("queryResolver:Credentials", &paginator.Params{
		First:  first,
		Last:   last,
		After:  after,
		Before: before,
	})
	err2.Check(err)

	res, err := r.db.GetCredentials(batch, agent.ID)
	err2.Check(err)

	return res.ToConnection(nil), nil
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

// CredentialConnection returns generated.CredentialConnectionResolver implementation.
func (r *Resolver) CredentialConnection() generated.CredentialConnectionResolver {
	return &credentialConnectionResolver{r}
}

// Event returns generated.EventResolver implementation.
func (r *Resolver) Event() generated.EventResolver { return &eventResolver{r} }

// Job returns generated.JobResolver implementation.
func (r *Resolver) Job() generated.JobResolver { return &jobResolver{r} }

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

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// Subscription returns generated.SubscriptionResolver implementation.
func (r *Resolver) Subscription() generated.SubscriptionResolver { return &subscriptionResolver{r} }

type basicMessageResolver struct{ *Resolver }
type credentialResolver struct{ *Resolver }
type credentialConnectionResolver struct{ *Resolver }
type eventResolver struct{ *Resolver }
type jobResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
type pairwiseResolver struct{ *Resolver }
type pairwiseConnectionResolver struct{ *Resolver }
type proofResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }
