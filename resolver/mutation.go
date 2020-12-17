package resolver

import (
	"context"

	db "github.com/findy-network/findy-agent-vault/db/model"
	"github.com/findy-network/findy-agent-vault/db/store"
	"github.com/findy-network/findy-agent-vault/graph/model"
	"github.com/findy-network/findy-agent-vault/utils"
	"github.com/lainio/err2"
)

func (r *mutationResolver) addEvent(job *db.Job, description string) (err error) {
	var connectionID, jobID *string
	if job != nil {
		connectionID = job.ConnectionID
		jobID = &job.ID
	}
	// TODO: event subscription
	_, err = r.db.AddEvent(&db.Event{
		Read:         false,
		Description:  description,
		ConnectionID: connectionID,
		JobID:        jobID,
	})
	return err
}

func (r *mutationResolver) addJob(job *db.Job, description string) (err error) {
	defer err2.Return(&err)
	job, err = r.db.AddJob(job)
	err2.Check(err)

	err2.Check(r.addEvent(job, description))

	return
}

func (r *mutationResolver) updateJob(job *db.Job, description string) (err error) {
	defer err2.Return(&err)
	job, err = r.db.UpdateJob(job)
	err2.Check(err)

	err2.Check(r.addEvent(job, description))

	return
}

func (r *mutationResolver) markEventRead(ctx context.Context, input model.MarkReadInput) (e *model.Event, err error) {
	defer err2.Return(&err)

	agent, err := store.GetAgent(ctx, r.db)
	err2.Check(err)

	utils.LogMed().Infof(
		"mutationResolver:MarkEventRead for tenant %s, event: %s",
		agent.ID,
		input.ID,
	)

	event, err := r.db.MarkEventRead(input.ID, agent.ID)
	err2.Check(err)

	return event.ToNode(), nil
}

func (r *mutationResolver) invite(ctx context.Context) (res *model.InvitationResponse, err error) {
	defer err2.Return(&err)
	utils.LogMed().Info("mutationResolver:Invite")

	agent, err := store.GetAgent(ctx, r.db)
	err2.Check(err)

	str, id, err := r.agency.Invite(ctx)
	err2.Check(err)

	img, err := utils.StrToQRCode(str)
	err2.Check(err)

	err2.Check(r.addJob(
		db.NewJob(id, agent.ID, &db.Job{
			ProtocolType:  model.ProtocolTypeConnection,
			InitiatedByUs: true,
			Status:        model.JobStatusWaiting,
			Result:        model.JobResultNone,
		}),
		"Created connection invitation",
	))

	res = &model.InvitationResponse{
		Invitation: str,
		ImageB64:   img,
	}

	return
}

func (r *mutationResolver) connect(ctx context.Context, input model.ConnectInput) (res *model.Response, err error) {
	defer err2.Return(&err)
	utils.LogMed().Info("mutationResolver:Connect")

	agent, err := store.GetAgent(ctx, r.db)
	err2.Check(err)

	id, err := r.agency.Connect(ctx, input.Invitation)
	err2.Check(err)

	err2.Check(r.addJob(
		db.NewJob(id, agent.ID, &db.Job{
			ProtocolType:  model.ProtocolTypeConnection,
			InitiatedByUs: false,
			Status:        model.JobStatusWaiting,
			Result:        model.JobResultNone,
		}),
		"Sent connection request",
	))

	res = &model.Response{Ok: true}
	return
}

func (r *mutationResolver) sendMessage(ctx context.Context, input model.MessageInput) (res *model.Response, err error) {
	defer err2.Return(&err)
	utils.LogMed().Info("mutationResolver:SendMessage")

	agent, err := store.GetAgent(ctx, r.db)
	err2.Check(err)

	id, err := r.agency.SendMessage(ctx, input.ConnectionID, input.Message)
	err2.Check(err)

	err2.Check(r.addJob(
		db.NewJob(id, agent.ID, &db.Job{
			ConnectionID:  &input.ConnectionID,
			ProtocolType:  model.ProtocolTypeBasicMessage,
			InitiatedByUs: true,
			Status:        model.JobStatusWaiting,
			Result:        model.JobResultNone,
		}),
		"Sent basic message",
	))

	res = &model.Response{Ok: true}
	return
}

func (r *mutationResolver) resume(ctx context.Context, input model.ResumeJobInput) (res *model.Response, err error) {
	defer err2.Return(&err)
	utils.LogMed().Info("mutationResolver:Resume")

	agent, err := store.GetAgent(ctx, r.db)
	err2.Check(err)

	job, err := r.db.GetJob(input.ID, agent.ID)
	err2.Check(err)

	desc := "Accepted"
	if !input.Accept {
		desc = "Declined"
	}

	switch job.ProtocolType {
	case model.ProtocolTypeCredential:
		err2.Check(r.agency.ResumeCredentialOffer(ctx, job.ID, input.Accept))
		desc += " credential"
	case model.ProtocolTypeProof:
		err2.Check(r.agency.ResumeProofRequest(ctx, job.ID, input.Accept))
		desc += " proof"
	case model.ProtocolTypeBasicMessage:
	case model.ProtocolTypeConnection:
	case model.ProtocolTypeNone:
		// N/A
		return
	}

	job.Status = model.JobStatusWaiting
	job.Result = model.JobResultNone

	err2.Check(r.updateJob(job, desc))

	res = &model.Response{Ok: true}

	return res, err
}

// ************* For testing: TODO: enable only with feature flag *************
func (r *mutationResolver) addRandomEvent(ctx context.Context) (ok bool, err error) {
	_, err = store.GetAgent(ctx, r.db)
	err2.Check(err)

	// TODO
	return
}

func (r *mutationResolver) addRandomMessage(ctx context.Context) (ok bool, err error) {
	_, err = store.GetAgent(ctx, r.db)
	err2.Check(err)

	// TODO
	return
}

func (r *mutationResolver) addRandomCredential(ctx context.Context) (ok bool, err error) {
	_, err = store.GetAgent(ctx, r.db)
	err2.Check(err)

	// TODO
	return
}

func (r *mutationResolver) addRandomProof(ctx context.Context) (ok bool, err error) {
	_, err = store.GetAgent(ctx, r.db)
	err2.Check(err)

	// TODO
	return
}
