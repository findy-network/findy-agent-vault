package model

import (
	"time"

	"github.com/findy-network/findy-agent-vault/graph/model"
	"github.com/findy-network/findy-agent-vault/paginator"
)

type Jobs struct {
	Jobs            []*Job
	HasNextPage     bool
	HasPreviousPage bool
}

type Job struct {
	*base
	ProtocolType  model.ProtocolType `faker:"oneof: NONE,NONE"`
	ProtocolID    *string            `faker:"-"`
	ConnectionID  *string            `faker:"-"`
	Status        model.JobStatus    `faker:"oneof: COMPLETE,COMPLETE"`
	Result        model.JobResult    `faker:"oneof: SUCCESS,SUCCESS"`
	InitiatedByUs bool
	Updated       time.Time
}

func NewJob(j *Job) *Job {
	if j != nil {
		return j.copy()
	}
	return &Job{base: &base{}}
}

func (j *Job) copy() (n *Job) {
	n = NewJob(nil)
	if j.base != nil {
		n.base = j.base.copy()
	}
	n.ProtocolType = j.ProtocolType
	if j.ProtocolID != nil {
		protocolID := *j.ProtocolID
		n.ProtocolID = &protocolID
	}
	if j.ConnectionID != nil {
		connectionID := *j.ConnectionID
		n.ConnectionID = &connectionID
	}
	n.Status = j.Status
	n.Result = j.Result
	n.InitiatedByUs = j.InitiatedByUs
	n.Updated = j.Updated
	return n
}

func (j *Job) ToEdge() *model.JobEdge {
	cursor := paginator.CreateCursor(j.Cursor, model.Job{})
	return &model.JobEdge{
		Cursor: cursor,
		Node:   j.ToNode(),
	}
}

func (j *Job) ToNode() *model.Job {
	return &model.Job{
		ID:        j.ID,
		Protocol:  j.ProtocolType,
		Status:    j.Status,
		Result:    j.Result,
		CreatedMs: timeToString(&j.Created),
		UpdatedMs: timeToString(&j.Updated),
	}
}

func (j *Jobs) ToConnection(id *string, completed *bool) *model.JobConnection {
	totalCount := len(j.Jobs)

	edges := make([]*model.JobEdge, totalCount)
	nodes := make([]*model.Job, totalCount)
	for index, event := range j.Jobs {
		edge := event.ToEdge()
		edges[index] = edge
		nodes[index] = edge.Node
	}

	var startCursor, endCursor *string
	if len(edges) > 0 {
		startCursor = &edges[0].Cursor
		endCursor = &edges[len(edges)-1].Cursor
	}
	return &model.JobConnection{
		ConnectionID: id,
		Completed:    completed,
		Edges:        edges,
		Nodes:        nodes,
		PageInfo: &model.PageInfo{
			EndCursor:       endCursor,
			HasNextPage:     j.HasNextPage,
			HasPreviousPage: j.HasPreviousPage,
			StartCursor:     startCursor,
		},
	}
}
