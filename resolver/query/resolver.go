package query

import (
	"context"
	"encoding/base64"

	"github.com/findy-network/findy-agent-vault/db/store"
	"github.com/findy-network/findy-agent-vault/graph/model"
	"github.com/findy-network/findy-agent-vault/paginator"
	"github.com/findy-network/findy-agent-vault/resolver/invitation"
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

func (r *Resolver) Connections(ctx context.Context, after, before *string, first, last *int) (c *model.PairwiseConnection, err error) {
	defer err2.Handle(&err, func() {})

	tenant := try.To1(r.GetAgent(ctx))

	utils.LogLow().Info("queryResolver:Connections for tenant: ", tenant.ID)

	batch := try.To1(paginator.Validate("queryResolver:Connections", &paginator.Params{
		First:  first,
		Last:   last,
		After:  after,
		Before: before,
		Object: model.Pairwise{},
	}))

	res := try.To1(r.db.GetConnections(batch, tenant.ID))

	return res.ToConnection(), nil
}

func (r *Resolver) Connection(ctx context.Context, id string) (c *model.Pairwise, err error) {
	defer err2.Handle(&err)

	tenant := try.To1(r.GetAgent(ctx))

	utils.LogLow().Infof("queryResolver:Connection id: %s for tenant %s", id, tenant.ID)

	conn := try.To1(r.db.GetConnection(id, tenant.ID))

	return conn.ToNode(), nil
}

func (r *Resolver) Credential(ctx context.Context, id string) (c *model.Credential, err error) {
	defer err2.Handle(&err)

	tenant := try.To1(r.GetAgent(ctx))

	utils.LogLow().Infof("queryResolver:Credential id: %s for tenant %s", id, tenant.ID)

	cred := try.To1(r.db.GetCredential(id, tenant.ID))

	return cred.ToNode(), nil
}

func (r *Resolver) Credentials(
	ctx context.Context,
	after, before *string,
	first, last *int,
) (c *model.CredentialConnection, err error) {
	defer err2.Handle(&err, func() {})

	tenant := try.To1(r.GetAgent(ctx))

	utils.LogLow().Info("queryResolver:Credentials for tenant: ", tenant.ID)

	batch := try.To1(paginator.Validate("queryResolver:Credentials", &paginator.Params{
		First:  first,
		Last:   last,
		After:  after,
		Before: before,
		Object: model.Credential{},
	}))

	res := try.To1(r.db.GetCredentials(batch, tenant.ID, nil))

	return res.ToConnection(nil), nil
}

func (r *Resolver) Proof(ctx context.Context, id string) (c *model.Proof, err error) {
	defer err2.Handle(&err)

	tenant := try.To1(r.GetAgent(ctx))

	utils.LogLow().Infof("queryResolver:Proof id: %s for tenant %s", id, tenant.ID)

	cred := try.To1(r.db.GetProof(id, tenant.ID))

	return cred.ToNode(), nil
}

func (r *Resolver) Message(ctx context.Context, id string) (c *model.BasicMessage, err error) {
	defer err2.Handle(&err)

	tenant := try.To1(r.GetAgent(ctx))

	utils.LogLow().Infof("queryResolver:Message id: %s for tenant %s", id, tenant.ID)

	msg := try.To1(r.db.GetMessage(id, tenant.ID))

	return msg.ToNode(), nil
}

func (r *Resolver) Events(ctx context.Context, after, before *string, first, last *int) (e *model.EventConnection, err error) {
	defer err2.Handle(&err, func() {})

	tenant := try.To1(r.GetAgent(ctx))

	utils.LogLow().Info("queryResolver:Events for tenant: ", tenant.ID)

	batch := try.To1(paginator.Validate("queryResolver:Events", &paginator.Params{
		First:  first,
		Last:   last,
		After:  after,
		Before: before,
		Object: model.Event{},
	}))

	res := try.To1(r.db.GetEvents(batch, tenant.ID, nil))

	return res.ToConnection(nil), nil
}

func (r *Resolver) Event(ctx context.Context, id string) (e *model.Event, err error) {
	defer err2.Handle(&err)

	tenant := try.To1(r.GetAgent(ctx))

	utils.LogLow().Infof("queryResolver:Event id: %s for tenant %s", id, tenant.ID)

	event := try.To1(r.db.GetEvent(id, tenant.ID))

	return event.ToNode(), nil
}

func (r *Resolver) Jobs(
	ctx context.Context,
	after, before *string,
	first, last *int,
	completed *bool,
) (e *model.JobConnection, err error) {
	defer err2.Handle(&err, func() {})

	tenant := try.To1(r.GetAgent(ctx))

	utils.LogLow().Info("queryResolver:Jobs for tenant: ", tenant.ID)

	batch := try.To1(paginator.Validate("queryResolver:Jobs", &paginator.Params{
		First:  first,
		Last:   last,
		After:  after,
		Before: before,
		Object: model.Job{},
	}))

	res := try.To1(r.db.GetJobs(batch, tenant.ID, nil, completed))

	return res.ToConnection(nil, completed), nil
}

func (r *Resolver) Job(ctx context.Context, id string) (e *model.Job, err error) {
	defer err2.Handle(&err)

	tenant := try.To1(r.GetAgent(ctx))

	utils.LogLow().Infof("queryResolver:Job id: %s for tenant %s", id, tenant.ID)

	job := try.To1(r.db.GetJob(id, tenant.ID))

	return job.ToNode(), nil
}

func (r *Resolver) User(ctx context.Context) (u *model.User, err error) {
	defer err2.Handle(&err)

	tenant := try.To1(r.GetAgent(ctx))

	utils.LogLow().Infof("queryResolver:User tenant %s", tenant.ID)

	return tenant.ToNode(), nil
}

func (r *Resolver) Endpoint(ctx context.Context, payload string) (i *model.InvitationResponse, err error) {
	defer err2.Handle(&err)

	tenant := try.To1(r.GetAgent(ctx))

	utils.LogLow().Infof("queryResolver:Endpoint tenant %s", tenant.ID)

	if decoded, err := base64.StdEncoding.DecodeString(payload); err == nil {
		payload = string(decoded)
	}
	return invitation.FromURLParam(payload)
}
