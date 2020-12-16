package mock

import (
	"errors"
	"time"

	"github.com/bxcodec/faker/v3"
	"github.com/findy-network/findy-agent-vault/db/model"
	"github.com/findy-network/findy-agent-vault/paginator"
)

type mockJob struct {
	*base
	job *model.Job
}

func (j *mockJob) Created() uint64 {
	return model.TimeToCursor(&j.job.Created)
}

func (j *mockJob) Identifier() string {
	return j.job.ID
}

func newJob(j *model.Job) *mockJob {
	var job *model.Job
	if j != nil {
		job = model.NewJob(j)
	}
	return &mockJob{base: &base{}, job: job}
}

func (j *mockJob) Copy() apiObject {
	return newJob(j.job)
}

func (j *mockJob) Job() *model.Job {
	return j.job
}

func (m *mockData) AddJob(j *model.Job) (*model.Job, error) {
	agent := m.agents[j.TenantID]

	n := model.NewJob(j)
	n.ID = faker.UUIDHyphenated()
	n.Created = time.Now().UTC()
	n.Cursor = model.TimeToCursor(&n.Created)
	agent.jobs.append(newJob(n))

	return n, nil
}

func (m *mockData) UpdateJob(j *model.Job) (*model.Job, error) {
	agent := m.agents[j.TenantID]

	object := agent.jobs.objectForID(j.ID)
	if object == nil {
		return nil, errors.New("not found job for id: " + j.ID)
	}
	var protocolID, connectionID *string
	if j.ProtocolID != nil {
		p := *j.ProtocolID
		protocolID = &p
	}
	if j.ConnectionID != nil {
		c := *j.ConnectionID
		connectionID = &c
	}
	updated := object.Copy()
	job := updated.Job()
	job.ProtocolID = protocolID
	job.ConnectionID = connectionID
	job.Status = j.Status
	job.Result = j.Result
	job.Updated = time.Now().UTC()

	if !agent.jobs.replaceObjectForID(j.ID, updated) {
		return nil, errors.New("not found job for id: " + j.ID)
	}
	return updated.Job(), nil
}

func (m *mockData) GetJob(id, tenantID string) (*model.Job, error) {
	agent := m.agents[tenantID]

	j := agent.jobs.objectForID(id)
	if j == nil {
		return nil, errors.New("not found job for id: " + id)
	}
	return j.Job(), nil
}

func (m *mockItems) getJobs(
	info *paginator.BatchInfo,
	filter func(item apiObject) bool,
) (connections *model.Jobs, err error) {
	state, hasNextPage, hasPreviousPage := m.jobs.getObjects(info, filter)
	res := make([]*model.Job, len(state.objects))
	for i := range state.objects {
		res[i] = state.objects[i].Copy().Job()
	}

	c := &model.Jobs{
		Jobs:            res,
		HasNextPage:     hasNextPage,
		HasPreviousPage: hasPreviousPage,
	}
	return c, nil
}

func jobConnectionFilter(id string) func(item apiObject) bool {
	return func(item apiObject) bool {
		j := item.Job()
		if j.ConnectionID != nil && *j.ConnectionID == id {
			return true
		}
		return false
	}
}

func (m *mockData) GetJobs(info *paginator.BatchInfo, tenantID string, connectionID *string) (connections *model.Jobs, err error) {
	agent := m.agents[tenantID]

	if connectionID == nil {
		return agent.getJobs(info, nil)
	}
	return agent.getJobs(info, jobConnectionFilter(*connectionID))
}

func (m *mockData) GetJobCount(tenantID string, connectionID *string) (int, error) {
	agent := m.agents[tenantID]

	if connectionID == nil {
		return agent.jobs.count(nil), nil
	}
	return agent.jobs.count(jobConnectionFilter(*connectionID)), nil
}
