package resolver

import (
	"context"

	"github.com/findy-network/findy-agent-vault/db/store"
	"github.com/findy-network/findy-agent-vault/graph/model"
	"github.com/findy-network/findy-agent-vault/paginator"
	"github.com/findy-network/findy-agent-vault/utils"
	"github.com/lainio/err2"
)

func (r *pairwiseResolver) credentials(
	ctx context.Context,
	obj *model.Pairwise,
	after, before *string,
	first, last *int,
) (c *model.CredentialConnection, err error) {
	defer err2.Return(&err)

	agent, err := store.GetAgent(ctx, r.db)
	err2.Check(err)

	utils.LogMed().Infof("pairwiseResolver:Credentials for tenant: %s, connection %s", agent.ID, obj.ID)

	batch, err := paginator.Validate("pairwiseResolver:Credentials", &paginator.Params{
		First:  first,
		Last:   last,
		After:  after,
		Before: before,
	})
	err2.Check(err)

	res, err := r.db.GetCredentials(batch, agent.ID, &obj.ID)
	err2.Check(err)

	return res.ToConnection(&obj.ID), nil
}

func (r *pairwiseResolver) events(
	ctx context.Context,
	obj *model.Pairwise,
	after, before *string,
	first, last *int,
) (e *model.EventConnection, err error) {
	defer err2.Return(&err)

	agent, err := store.GetAgent(ctx, r.db)
	err2.Check(err)

	utils.LogMed().Infof("pairwiseResolver:Events for tenant: %s, connection %s", agent.ID, obj.ID)

	batch, err := paginator.Validate("pairwiseResolver:Events", &paginator.Params{
		First:  first,
		Last:   last,
		After:  after,
		Before: before,
	})
	err2.Check(err)

	res, err := r.db.GetEvents(batch, agent.ID, &obj.ID)
	err2.Check(err)

	return res.ToConnection(&obj.ID), nil
}

func (r *pairwiseResolver) messages(
	ctx context.Context,
	obj *model.Pairwise,
	after, before *string,
	first, last *int,
) (e *model.BasicMessageConnection, err error) {
	defer err2.Return(&err)

	agent, err := store.GetAgent(ctx, r.db)
	err2.Check(err)

	utils.LogMed().Infof("pairwiseResolver:Messages for tenant: %s, connection %s", agent.ID, obj.ID)

	batch, err := paginator.Validate("pairwiseResolver:Messages", &paginator.Params{
		First:  first,
		Last:   last,
		After:  after,
		Before: before,
	})
	err2.Check(err)

	res, err := r.db.GetMessages(batch, agent.ID, &obj.ID)
	err2.Check(err)

	return res.ToConnection(&obj.ID), nil
}
