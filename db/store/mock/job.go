package mock

import (
	"errors"
	"time"

	"github.com/findy-network/findy-agent-vault/db/model"
	"github.com/findy-network/findy-agent-vault/db/store"
	graph "github.com/findy-network/findy-agent-vault/graph/model"
	"github.com/findy-network/findy-agent-vault/paginator"
	"github.com/findy-network/findy-agent-vault/utils"
	"github.com/lainio/err2"
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
		job = model.NewJob(j.ID, j.TenantID, j)
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
	agent := m.agents.get(j.TenantID)

	n := model.NewJob(j.ID, j.TenantID, j)
	n.Created = time.Now().UTC()
	n.Cursor = model.TimeToCursor(&n.Created)
	agent.jobs.append(newJob(n))

	return n, nil
}

func (m *mockData) UpdateJob(arg *model.Job) (*model.Job, error) {
	agent := m.agents.get(arg.TenantID)

	object := agent.jobs.objectForID(arg.ID)
	if object == nil {
		return nil, store.NewError(store.ErrCodeNotFound, "not found job for id: "+arg.ID)
	}
	updated := object.Copy()
	job := updated.Job()
	job.ProtocolConnectionID = utils.CopyStrPtr(arg.ProtocolConnectionID)
	job.ProtocolCredentialID = utils.CopyStrPtr(arg.ProtocolCredentialID)
	job.ProtocolProofID = utils.CopyStrPtr(arg.ProtocolProofID)
	job.ProtocolMessageID = utils.CopyStrPtr(arg.ProtocolMessageID)
	job.ConnectionID = utils.CopyStrPtr(arg.ConnectionID)
	job.Status = arg.Status
	job.Result = arg.Result
	job.Updated = time.Now().UTC()

	if !agent.jobs.replaceObjectForID(arg.ID, updated) {
		panic("not found job for id: " + arg.ID)
	}
	return updated.Job(), nil
}

func (m *mockData) GetJob(id, tenantID string) (*model.Job, error) {
	agent := m.agents.get(tenantID)

	j := agent.jobs.objectForID(id)
	if j == nil {
		return nil, store.NewError(store.ErrCodeNotFound, "not found job for id: "+id)
	}
	return j.Job(), nil
}

func (m *mockItems) getJobs(
	info *paginator.BatchInfo,
	filter func(item apiObject) bool,
) (jobs *model.Jobs, err error) {
	state, hasNextPage, hasPreviousPage := m.jobs.getObjects(info, filter)
	res := make([]*model.Job, len(state.objects))
	for i := range state.objects {
		res[i] = state.objects[i].Copy().Job()
	}

	jobs = &model.Jobs{
		Jobs:            res,
		HasNextPage:     hasNextPage,
		HasPreviousPage: hasPreviousPage,
	}
	return
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
	agent := m.agents.get(tenantID)

	if connectionID == nil {
		return agent.getJobs(info, jobFilter(completed))
	}
	return agent.getJobs(info, jobConnectionFilter(*connectionID, completed))
}

func (m *mockData) GetJobCount(tenantID string, connectionID *string, completed *bool) (int, error) {
	agent := m.agents.get(tenantID)

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

func (m *mockData) GetJobOutput(id, tenantID string, protocolType graph.ProtocolType) (output *model.JobOutput, err error) {
	defer err2.Return(&err)

	job, err := m.GetJob(id, tenantID)
	if err != nil {
		return nil, err
	}
	switch protocolType {
	case graph.ProtocolTypeConnection:
		connection, err := m.GetConnection(*job.ProtocolConnectionID, tenantID)
		err2.Check(err)
		return &model.JobOutput{Connection: connection}, nil
	case graph.ProtocolTypeCredential:
		credential, err := m.GetCredential(*job.ProtocolCredentialID, tenantID)
		err2.Check(err)
		return &model.JobOutput{Credential: credential}, nil
	case graph.ProtocolTypeProof:
		proof, err := m.GetProof(*job.ProtocolProofID, tenantID)
		err2.Check(err)
		return &model.JobOutput{Proof: proof}, nil
	case graph.ProtocolTypeBasicMessage:
		message, err := m.GetMessage(*job.ProtocolMessageID, tenantID)
		err2.Check(err)
		return &model.JobOutput{Message: message}, nil
	case graph.ProtocolTypeNone:
		break
	}
	return &model.JobOutput{}, nil
}
