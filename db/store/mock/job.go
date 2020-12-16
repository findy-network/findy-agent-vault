package mock

import (
	"errors"
	"time"

	"github.com/bxcodec/faker/v3"
	"github.com/findy-network/findy-agent-vault/db/model"
	graph "github.com/findy-network/findy-agent-vault/graph/model"
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

func jobFilter(completed *bool) func(item apiObject) bool {
	fetchAll := completed != nil && *completed
	return func(item apiObject) bool {
		j := item.Job()
		if !fetchAll {
			return j.Status != graph.JobStatusComplete
		}
		return true
	}
}

func jobConnectionFilter(id string, completed *bool) func(item apiObject) bool {
	fetchAll := completed != nil && *completed
	return func(item apiObject) bool {
		j := item.Job()
		if j.ConnectionID != nil && *j.ConnectionID == id {
			if !fetchAll {
				return j.Status != graph.JobStatusComplete
			}
			return true
		}
		return false
	}
}

func (m *mockData) GetJobs(
	info *paginator.BatchInfo,
	tenantID string,
	connectionID *string,
	completed *bool,
) (connections *model.Jobs, err error) {
	agent := m.agents[tenantID]

	if connectionID == nil {
		return agent.getJobs(info, jobFilter(completed))
	}
	return agent.getJobs(info, jobConnectionFilter(*connectionID, completed))
}

func (m *mockData) GetJobCount(tenantID string, connectionID *string, completed *bool) (int, error) {
	agent := m.agents[tenantID]

	if connectionID == nil {
		return agent.jobs.count(jobFilter(completed)), nil
	}
	return agent.jobs.count(jobConnectionFilter(*connectionID, completed)), nil
}

func (m *mockData) GetConnectionForJob(id, tenantID string) (*model.Connection, error) {
	job, err := m.GetJob(id, tenantID)
	if err != nil {
		return nil, err
	}
	if job.ConnectionID != nil {
		return m.GetConnection(*job.ConnectionID, tenantID)
	}
	return nil, errors.New("no connection found for job id: " + id)
}
