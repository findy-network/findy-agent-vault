package model

import (
	"strconv"

	"github.com/findy-network/findy-agent-vault/graph/model"
)

type InternalJob struct {
	ID            string
	ProtocolType  model.ProtocolType
	Status        model.JobStatus
	Result        model.JobResult
	Details       *model.JobDetails
	InitiatedByUs bool
	CreatedMs     int64
	UpdatedMs     int64
}

func (j *InternalJob) Created() int64 {
	return j.CreatedMs
}

func (j *InternalJob) Identifier() string {
	return j.ID
}

func (j *InternalJob) Pairwise() *InternalPairwise {
	panic("Job is not pairwise")
}

func (j *InternalJob) Event() *InternalEvent {
	panic("Job is not event")
}

func (j *InternalJob) Job() *InternalJob {
	return j
}

func (j *InternalJob) Copy() *InternalJob {
	newJob := &InternalJob{
		ID:            j.ID,
		ProtocolType:  j.ProtocolType,
		Status:        j.Status,
		Result:        j.Result,
		Details:       j.copyDetails(),
		InitiatedByUs: j.InitiatedByUs,
		CreatedMs:     j.CreatedMs,
		UpdatedMs:     j.UpdatedMs,
	}
	return newJob
}

func (j *InternalJob) copyDetails() *model.JobDetails {
	newDetails := &model.JobDetails{
		PairwiseID:       j.Details.PairwiseID,
		CredDefID:        j.Details.CredDefID,
		CredentialValues: make([]*model.CredentialValue, 0),
		Verified:         j.Details.Verified,
	}
	for _, v := range j.Details.CredentialValues {
		newDetails.CredentialValues = append(newDetails.CredentialValues, &model.CredentialValue{
			Name:  v.Name,
			Value: v.Value,
		})
	}
	return newDetails
}

func (j *InternalJob) ToEdge() *model.JobEdge {
	cursor := CreateCursor(j.CreatedMs, model.Job{})
	return &model.JobEdge{
		Cursor: cursor,
		Node:   j.ToNode(),
	}
}

func (j *InternalJob) ToNode() *model.Job {
	createdStr := strconv.FormatInt(j.CreatedMs, 10)
	updatedStr := strconv.FormatInt(j.UpdatedMs, 10)
	return &model.Job{
		ID:            j.ID,
		Protocol:      j.ProtocolType,
		Status:        j.Status,
		Result:        j.Result,
		Details:       j.copyDetails(),
		InitiatedByUs: j.InitiatedByUs,
		CreatedMs:     createdStr,
		UpdatedMs:     updatedStr,
	}
}

func (i *Items) JobForID(id string) (node *model.Job) {
	i.mutex.RLock()
	defer i.mutex.RUnlock()

	for _, item := range i.items {
		if item.Identifier() == id {
			node = item.Job().ToNode()
			break
		}
	}
	return
}

func (i *Items) JobConnection(after, before int) *model.JobConnection {
	i.mutex.RLock()
	result := i.items[after:before]
	totalCount := len(result)

	edges := make([]*model.JobEdge, totalCount)
	nodes := make([]*model.Job, totalCount)
	for index, job := range result {
		node := job.Job().ToNode()
		edges[index] = &model.JobEdge{
			Cursor: CreateCursor(job.Job().CreatedMs, model.Job{}),
			Node:   node,
		}
		nodes[index] = node
	}
	i.mutex.RUnlock()

	var startCursor, endCursor *string
	if totalCount > 0 {
		startCursor = &edges[0].Cursor
		endCursor = &edges[totalCount-1].Cursor
	}
	c := &model.JobConnection{
		Edges: edges,
		Nodes: nodes,
		PageInfo: &model.PageInfo{
			EndCursor:       endCursor,
			HasNextPage:     edges[len(edges)-1].Node.ID != i.LastID(),
			HasPreviousPage: edges[0].Node.ID != i.FirstID(),
			StartCursor:     startCursor,
		},
		TotalCount: totalCount,
	}
	return c
}
