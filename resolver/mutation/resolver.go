package mutation

import (
	"context"

	agency "github.com/findy-network/findy-agent-vault/agency/model"
	dbModel "github.com/findy-network/findy-agent-vault/db/model"
	"github.com/findy-network/findy-agent-vault/db/store"
	"github.com/findy-network/findy-agent-vault/graph/model"
	"github.com/findy-network/findy-agent-vault/resolver/invitation"
	"github.com/findy-network/findy-agent-vault/resolver/query/agent"
	"github.com/findy-network/findy-agent-vault/resolver/update"
	"github.com/findy-network/findy-agent-vault/utils"
	"github.com/lainio/err2"
	"github.com/lainio/err2/try"
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
	defer err2.Handle(&err)

	tenant := try.To1(r.GetAgent(ctx))

	utils.LogLow().Infof(
		"mutationResolver:MarkEventRead for tenant %s, event: %s",
		tenant.ID,
		input.ID,
	)

	event := try.To1(r.db.MarkEventRead(input.ID, tenant.ID))

	return event.ToNode(), nil
}

func (r *Resolver) Invite(ctx context.Context) (res *model.InvitationResponse, err error) {
	defer err2.Handle(&err)
	utils.LogLow().Info("mutationResolver:Invite")

	tenant := try.To1(r.GetAgent(ctx))

	data := try.To1(r.agency.Invite(r.AgencyAuth(tenant)))

	res = try.To1(invitation.FromAgency(data))

	try.To(r.AddJob(
		&dbModel.Job{
			Base:          dbModel.Base{ID: data.ID, TenantID: tenant.ID},
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
	defer err2.Handle(&err)
	utils.LogLow().Info("mutationResolver:Connect")

	tenant := try.To1(r.GetAgent(ctx))

	id := try.To1(r.agency.Connect(r.AgencyAuth(tenant), input.Invitation))

	try.To(r.AddJob(
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
	defer err2.Handle(&err)
	utils.LogLow().Info("mutationResolver:SendMessage")

	tenant := try.To1(r.GetAgent(ctx))

	try.To1(r.agency.SendMessage(r.AgencyAuth(tenant), input.ConnectionID, input.Message))

	res = &model.Response{Ok: true}
	return
}

func (r *Resolver) Resume(ctx context.Context, input model.ResumeJobInput) (res *model.Response, err error) {
	defer err2.Handle(&err)
	utils.LogLow().Info("mutationResolver:Resume")

	tenant := try.To1(r.GetAgent(ctx))

	job := try.To1(r.db.GetJob(input.ID, tenant.ID))

	jobInfo := &agency.JobInfo{
		TenantID:     tenant.ID,
		JobID:        job.ID,
		ConnectionID: *job.ConnectionID,
	}

	switch job.ProtocolType {
	case model.ProtocolTypeCredential:
		try.To(r.agency.ResumeCredentialOffer(r.AgencyAuth(tenant), jobInfo, input.Accept))
	case model.ProtocolTypeProof:
		try.To(r.agency.ResumeProofRequest(r.AgencyAuth(tenant), jobInfo, input.Accept))
	case model.ProtocolTypeBasicMessage:
	case model.ProtocolTypeConnection:
	case model.ProtocolTypeNone:
		// N/A
		return
	}

	res = &model.Response{Ok: true}

	return res, err
}
