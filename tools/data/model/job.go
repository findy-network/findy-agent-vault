package model

import (
	"strconv"

	"github.com/findy-network/findy-agent-vault/graph/model"
	"github.com/findy-network/findy-agent-vault/tools/utils"
)

type InternalJob struct {
	*BaseObject
	ProtocolType  model.ProtocolType
	ProtocolID    *string
	PairwiseID    *string
	Status        model.JobStatus
	Result        model.JobResult
	InitiatedByUs bool
	UpdatedMs     int64
}

func (j *InternalJob) Job() *InternalJob {
	return j
}

func (j *InternalJob) Copy() *InternalJob {
	newJob := &InternalJob{
		BaseObject: &BaseObject{
			ID:        j.ID,
			CreatedMs: j.CreatedMs,
		},
		ProtocolType:  j.ProtocolType,
		ProtocolID:    j.ProtocolID,
		PairwiseID:    j.PairwiseID,
		Status:        j.Status,
		Result:        j.Result,
		InitiatedByUs: j.InitiatedByUs,
		UpdatedMs:     j.UpdatedMs,
	}
	return newJob
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
		ID:        j.ID,
		Protocol:  j.ProtocolType,
		Status:    j.Status,
		Result:    j.Result,
		CreatedMs: createdStr,
		UpdatedMs: updatedStr,
	}
}

func (i *Items) IsJobInitiatedByUs(id string) (is *bool) {
	i.mutex.RLock()
	defer i.mutex.RUnlock()

	for _, item := range i.items {
		if item.Identifier() == id {
			jobInitiated := item.Job().InitiatedByUs
			is = &jobInitiated
			break
		}
	}

	return
}

func (i *Items) JobForID(id string) (edge *model.JobEdge) {
	i.mutex.RLock()
	defer i.mutex.RUnlock()

	for _, item := range i.items {
		if item.Identifier() == id {
			edge = item.Job().ToEdge()
			break
		}
	}
	return
}

func (i *Items) JobProtocolForID(id string) (t model.ProtocolType, protocolID *string) {
	i.mutex.RLock()
	defer i.mutex.RUnlock()

	t = model.ProtocolTypeNone
	var pID string

	for _, item := range i.items {
		if item.Identifier() == id {
			job := item.Job()
			t = job.ProtocolType
			if job.ProtocolID != nil {
				pID = *job.ProtocolID
				protocolID = &pID
			}
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
	var hasNextPage, hasPreviousPage bool
	if totalCount > 0 {
		startCursor = &edges[0].Cursor
		endCursor = &edges[totalCount-1].Cursor
		hasNextPage = edges[len(edges)-1].Node.ID != i.LastID()
		hasPreviousPage = edges[0].Node.ID != i.FirstID()
	}

	c := &model.JobConnection{
		Edges: edges,
		Nodes: nodes,
		PageInfo: &model.PageInfo{
			EndCursor:       endCursor,
			HasNextPage:     hasNextPage,
			HasPreviousPage: hasPreviousPage,
			StartCursor:     startCursor,
		},
		TotalCount: totalCount,
	}
	return c
}

func (i *Items) UpdateJob(id string, protocolID, pairwiseID *string, status model.JobStatus, result model.JobResult) bool {
	i.mutex.Lock()
	defer i.mutex.Unlock()

	for _, item := range i.items {
		if item.Identifier() != id {
			continue
		}
		job := item.Job()
		job.UpdatedMs = utils.CurrentTimeMs()
		job.Status = status
		job.Result = result
		job.ProtocolID = protocolID
		job.PairwiseID = pairwiseID
		return true
	}

	return false
}
