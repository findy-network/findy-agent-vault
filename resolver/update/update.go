package update

import (
	"github.com/findy-network/findy-agent-vault/db/model"
	"github.com/findy-network/findy-agent-vault/db/store"
	"github.com/findy-network/findy-agent-vault/resolver/query/agent"
	"github.com/findy-network/findy-agent-vault/utils"
	"github.com/lainio/err2"
	"github.com/lainio/err2/try"
)

type Updater struct {
	db               store.DB
	eventSubscribers *subscriberRegister
	*agent.Resolver
}

func NewUpdater(db store.DB, agentResolver *agent.Resolver) *Updater {
	return &Updater{
		db,
		newSubscriberRegister(),
		agentResolver,
	}
}

func (r *Updater) AddEvent(tenantID string, job *model.Job, description string) (err error) {
	defer err2.Handle(&err)
	var connectionID, jobID *string
	if job != nil {
		connectionID = job.ConnectionID
		jobID = &job.ID
	}
	event := try.To1(r.db.AddEvent(&model.Event{
		Base:         model.Base{TenantID: tenantID},
		Read:         false,
		Description:  description,
		ConnectionID: connectionID,
		JobID:        jobID,
	}))

	r.eventSubscribers.notify(tenantID, event)
	return err
}

func (r *Updater) AddJob(inputJob *model.Job, description string) (job *model.Job, err error) {
	defer err2.Handle(&err)

	utils.LogMed().Infof("Add job with ID %s for tenant %s", inputJob.ID, inputJob.TenantID)

	job = try.To1(r.db.AddJob(inputJob))

	try.To(r.AddEvent(job.TenantID, job, description))

	return
}

func (r *Updater) UpdateJob(job *model.Job, description string) (err error) {
	defer err2.Handle(&err)

	utils.LogMed().Infof("Update job with ID %s for tenant %s", job.ID, job.TenantID)

	job = try.To1(r.db.UpdateJob(job))

	try.To(r.AddEvent(job.TenantID, job, description))

	return
}
