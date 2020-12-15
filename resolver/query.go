package resolver

import (
	"context"

	"github.com/findy-network/findy-agent-vault/db/store"
	"github.com/findy-network/findy-agent-vault/graph/model"
	"github.com/findy-network/findy-agent-vault/paginator"
	"github.com/findy-network/findy-agent-vault/utils"
	"github.com/lainio/err2"
)

func (r *queryResolver) connections(ctx context.Context, after, before *string, first, last *int) (c *model.PairwiseConnection, err error) {
	defer err2.Return(&err)

	agent, err := store.GetAgent(ctx, r.db)
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

func (r *queryResolver) connection(ctx context.Context, id string) (c *model.Pairwise, err error) {
	defer err2.Return(&err)

	agent, err := store.GetAgent(ctx, r.db)
	err2.Check(err)

	utils.LogMed().Infof("queryResolver:Connection id: %s for tenant %s", id, agent.ID)

	conn, err := r.db.GetConnection(id, agent.ID)
	err2.Check(err)

	return conn.ToNode(), nil
}

func (r *queryResolver) credential(ctx context.Context, id string) (c *model.Credential, err error) {
	defer err2.Return(&err)

	agent, err := store.GetAgent(ctx, r.db)
	err2.Check(err)

	utils.LogMed().Infof("queryResolver:Credential id: %s for tenant %s", id, agent.ID)

	cred, err := r.db.GetCredential(id, agent.ID)
	err2.Check(err)

	return cred.ToNode(), nil
}

func (r *queryResolver) credentials(
	ctx context.Context,
	after, before *string,
	first, last *int,
) (c *model.CredentialConnection, err error) {
	defer err2.Return(&err)

	agent, err := store.GetAgent(ctx, r.db)
	err2.Check(err)

	utils.LogMed().Info("queryResolver:Credentials for tenant: ", agent.ID)

	batch, err := paginator.Validate("queryResolver:Credentials", &paginator.Params{
		First:  first,
		Last:   last,
		After:  after,
		Before: before,
	})
	err2.Check(err)

	res, err := r.db.GetCredentials(batch, agent.ID, nil)
	err2.Check(err)

	return res.ToConnection(nil), nil
}

func (r *queryResolver) events(ctx context.Context, after, before *string, first, last *int) (e *model.EventConnection, err error) {
	defer err2.Return(&err)

	agent, err := store.GetAgent(ctx, r.db)
	err2.Check(err)

	utils.LogMed().Info("queryResolver:Events for tenant: ", agent.ID)

	batch, err := paginator.Validate("queryResolver:Events", &paginator.Params{
		First:  first,
		Last:   last,
		After:  after,
		Before: before,
	})
	err2.Check(err)

	res, err := r.db.GetEvents(batch, agent.ID, nil)
	err2.Check(err)

	return res.ToConnection(nil), nil
}

func (r *queryResolver) event(ctx context.Context, id string) (e *model.Event, err error) {
	defer err2.Return(&err)

	agent, err := store.GetAgent(ctx, r.db)
	err2.Check(err)

	utils.LogMed().Infof("queryResolver:Event id: %s for tenant %s", id, agent.ID)

	event, err := r.db.GetEvent(id, agent.ID)
	err2.Check(err)

	return event.ToNode(), nil
}
