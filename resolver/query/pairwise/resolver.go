package pairwise

import (
	"context"

	"github.com/findy-network/findy-agent-vault/db/store"
	"github.com/findy-network/findy-agent-vault/graph/model"
	"github.com/findy-network/findy-agent-vault/paginator"
	"github.com/findy-network/findy-agent-vault/resolver/query/agent"
	"github.com/findy-network/findy-agent-vault/utils"
	"github.com/lainio/err2"
	"github.com/lainio/err2/try"
)

type Resolver struct {
	db store.DB
	*agent.Resolver
}

func NewResolver(db store.DB, agentResolver *agent.Resolver) *Resolver {
	return &Resolver{db, agentResolver}
}

func (r *Resolver) Credentials(
	ctx context.Context,
	obj *model.Pairwise,
	after, before *string,
	first, last *int,
) (c *model.CredentialConnection, err error) {
	defer err2.Return(&err)

	tenant := try.To1(r.GetAgent(ctx))

	utils.LogLow().Infof("pairwiseResolver:Credentials for tenant: %s, connection %s", tenant.ID, obj.ID)

	batch := try.To1(paginator.Validate("pairwiseResolver:Credentials", &paginator.Params{
		First:  first,
		Last:   last,
		After:  after,
		Before: before,
		Object: model.Credential{},
	}))

	res := try.To1(r.db.GetCredentials(batch, tenant.ID, &obj.ID))

	return res.ToConnection(&obj.ID), nil
}

func (r *Resolver) Proofs(
	ctx context.Context,
	obj *model.Pairwise,
	after, before *string,
	first, last *int,
) (c *model.ProofConnection, err error) {
	defer err2.Return(&err)

	tenant := try.To1(r.GetAgent(ctx))

	utils.LogLow().Infof("pairwiseResolver:Proofs for tenant: %s, connection %s", tenant.ID, obj.ID)

	batch := try.To1(paginator.Validate("pairwiseResolver:Proofs", &paginator.Params{
		First:  first,
		Last:   last,
		After:  after,
		Before: before,
		Object: model.Proof{},
	}))

	res := try.To1(r.db.GetProofs(batch, tenant.ID, &obj.ID))

	return res.ToConnection(&obj.ID), nil
}

func (r *Resolver) Messages(
	ctx context.Context,
	obj *model.Pairwise,
	after, before *string,
	first, last *int,
) (e *model.BasicMessageConnection, err error) {
	defer err2.Return(&err)

	tenant := try.To1(r.GetAgent(ctx))

	utils.LogLow().Infof("pairwiseResolver:Messages for tenant: %s, connection %s", tenant.ID, obj.ID)

	batch := try.To1(paginator.Validate("pairwiseResolver:Messages", &paginator.Params{
		First:  first,
		Last:   last,
		After:  after,
		Before: before,
		Object: model.BasicMessage{},
	}))

	res := try.To1(r.db.GetMessages(batch, tenant.ID, &obj.ID))

	return res.ToConnection(&obj.ID), nil
}

func (r *Resolver) Events(
	ctx context.Context,
	obj *model.Pairwise,
	after, before *string,
	first, last *int,
) (e *model.EventConnection, err error) {
	defer err2.Return(&err)

	tenant := try.To1(r.GetAgent(ctx))

	utils.LogLow().Infof("pairwiseResolver:Events for tenant: %s, connection %s", tenant.ID, obj.ID)

	batch := try.To1(paginator.Validate("pairwiseResolver:Events", &paginator.Params{
		First:  first,
		Last:   last,
		After:  after,
		Before: before,
		Object: model.Event{},
	}))

	res := try.To1(r.db.GetEvents(batch, tenant.ID, &obj.ID))

	return res.ToConnection(&obj.ID), nil
}

func (r *Resolver) Jobs(
	ctx context.Context,
	obj *model.Pairwise,
	after, before *string,
	first, last *int,
	completed *bool,
) (e *model.JobConnection, err error) {
	defer err2.Return(&err)

	tenant := try.To1(r.GetAgent(ctx))

	utils.LogLow().Infof("pairwiseResolver:Jobs for tenant: %s, connection %s", tenant.ID, obj.ID)

	batch := try.To1(paginator.Validate("pairwiseResolver:Jobs", &paginator.Params{
		First:  first,
		Last:   last,
		After:  after,
		Before: before,
		Object: model.Job{},
	}))

	res := try.To1(r.db.GetJobs(batch, tenant.ID, &obj.ID, completed))

	return res.ToConnection(&obj.ID, completed), nil
}
