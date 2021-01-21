package update

import (
	"github.com/findy-network/findy-agent-vault/db/model"
	"github.com/findy-network/findy-agent-vault/db/store"
	"github.com/findy-network/findy-agent-vault/resolver/agent"
	"github.com/findy-network/findy-agent-vault/utils"
	"github.com/lainio/err2"
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
	defer err2.Return(&err)
	var connectionID, jobID *string
	if job != nil {
		connectionID = job.ConnectionID
		jobID = &job.ID
	}
	event, err := r.db.AddEvent(model.NewEvent(tenantID, &model.Event{
		Read:         false,
		Description:  description,
		ConnectionID: connectionID,
		JobID:        jobID,
	}))
	err2.Check(err)

	r.eventSubscribers.notify(tenantID, event)
	return err
}

func (r *Updater) AddJob(job *model.Job, description string) (err error) {
	defer err2.Return(&err)

	utils.LogMed().Infof("Add job with ID %s for tenant %s", job.ID, job.TenantID)

	job, err = r.db.AddJob(job)
	err2.Check(err)

	err2.Check(r.AddEvent(job.TenantID, job, description))

	return
}

func (r *Updater) UpdateJob(job *model.Job, description string) (err error) {
	defer err2.Return(&err)

	utils.LogMed().Infof("Update job with ID %s for tenant %s", job.ID, job.TenantID)

	job, err = r.db.UpdateJob(job)
	err2.Check(err)

	err2.Check(r.AddEvent(job.TenantID, job, description))

	return
}
