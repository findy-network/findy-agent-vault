package mutation

import (
	"context"

	agency "github.com/findy-network/findy-agent-vault/agency/model"
	dbModel "github.com/findy-network/findy-agent-vault/db/model"
	"github.com/findy-network/findy-agent-vault/db/store"
	"github.com/findy-network/findy-agent-vault/graph/model"
	"github.com/findy-network/findy-agent-vault/resolver/query/agent"
	"github.com/findy-network/findy-agent-vault/resolver/update"
	"github.com/findy-network/findy-agent-vault/utils"
	"github.com/lainio/err2"
)

type Resolver struct {
	db     store.DB
	agency agency.Agency
	*agent.Resolver
	*update.Updater
}

func NewResolver(
	db store.DB,
	agencyInstance agency.Agency,
	agentResolver *agent.Resolver,
	updater *update.Updater,
) *Resolver {
	return &Resolver{db, agencyInstance, agentResolver, updater}
}

func (r *Resolver) MarkEventRead(ctx context.Context, input model.MarkReadInput) (e *model.Event, err error) {
	defer err2.Return(&err)

	tenant, err := r.GetAgent(ctx)
	err2.Check(err)

	utils.LogLow().Infof(
		"mutationResolver:MarkEventRead for tenant %s, event: %s",
		tenant.ID,
		input.ID,
	)

	event, err := r.db.MarkEventRead(input.ID, tenant.ID)
	err2.Check(err)

	return event.ToNode(), nil
}

func (r *Resolver) Invite(ctx context.Context) (res *model.InvitationResponse, err error) {
	defer err2.Return(&err)
	utils.LogLow().Info("mutationResolver:Invite")

	tenant, err := r.GetAgent(ctx)
	err2.Check(err)

	str, id, err := r.agency.Invite(r.AgencyAuth(tenant))
	err2.Check(err)

	res, err = utils.FromAriesInvitation(str)
	err2.Check(err)

	err2.Check(r.AddJob(
		&dbModel.Job{
			Base:          dbModel.Base{ID: id, TenantID: tenant.ID},
			ProtocolType:  model.ProtocolTypeConnection,
			InitiatedByUs: true,
			Status:        model.JobStatusWaiting,
			Result:        model.JobResultNone,
		},
		"Created connection invitation",
	))

	return
}

func (r *Resolver) Connect(ctx context.Context, input model.ConnectInput) (res *model.Response, err error) {
	defer err2.Return(&err)
	utils.LogLow().Info("mutationResolver:Connect")

	tenant, err := r.GetAgent(ctx)
	err2.Check(err)

	id, err := r.agency.Connect(r.AgencyAuth(tenant), input.Invitation)
	err2.Check(err)

	err2.Check(r.AddJob(
		&dbModel.Job{
			Base:          dbModel.Base{ID: id, TenantID: tenant.ID},
			ProtocolType:  model.ProtocolTypeConnection,
			InitiatedByUs: false,
			Status:        model.JobStatusWaiting,
			Result:        model.JobResultNone,
		},
		"Sent connection request",
	))

	res = &model.Response{Ok: true}
	return
}

func (r *Resolver) SendMessage(ctx context.Context, input model.MessageInput) (res *model.Response, err error) {
	defer err2.Return(&err)
	utils.LogLow().Info("mutationResolver:SendMessage")

	tenant, err := r.GetAgent(ctx)
	err2.Check(err)

	_, err = r.agency.SendMessage(r.AgencyAuth(tenant), input.ConnectionID, input.Message)
	err2.Check(err)

	res = &model.Response{Ok: true}
	return
}

func (r *Resolver) Resume(ctx context.Context, input model.ResumeJobInput) (res *model.Response, err error) {
	defer err2.Return(&err)
	utils.LogLow().Info("mutationResolver:Resume")

	tenant, err := r.GetAgent(ctx)
	err2.Check(err)

	job, err := r.db.GetJob(input.ID, tenant.ID)
	err2.Check(err)

	jobInfo := &agency.JobInfo{
		TenantID:     tenant.ID,
		JobID:        job.ID,
		ConnectionID: *job.ConnectionID,
	}

	switch job.ProtocolType {
	case model.ProtocolTypeCredential:
		err2.Check(r.agency.ResumeCredentialOffer(r.AgencyAuth(tenant), jobInfo, input.Accept))
	case model.ProtocolTypeProof:
		err2.Check(r.agency.ResumeProofRequest(r.AgencyAuth(tenant), jobInfo, input.Accept))
	case model.ProtocolTypeBasicMessage:
	case model.ProtocolTypeConnection:
	case model.ProtocolTypeNone:
		// N/A
		return
	}

	res = &model.Response{Ok: true}

	return res, err
}
