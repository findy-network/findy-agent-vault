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
	*Base
	ProtocolType         model.ProtocolType `faker:"oneof: NONE,NONE"`
	ProtocolConnectionID *string            `faker:"-"`
	ProtocolCredentialID *string            `faker:"-"`
	ProtocolProofID      *string            `faker:"-"`
	ProtocolMessageID    *string            `faker:"-"`
	ConnectionID         *string            `faker:"-"`
	Status               model.JobStatus    `faker:"oneof: COMPLETE,COMPLETE"`
	Result               model.JobResult    `faker:"oneof: SUCCESS,SUCCESS"`
	InitiatedByUs        bool
	Updated              time.Time
}

type JobOutput struct {
	Connection *Connection
	Credential *Credential
	Proof      *Proof
	Message    *Message
}

func NewJob(id, tenantID string, j *Job) *Job {
	defaultBase := &Base{ID: id, TenantID: tenantID}
	if j != nil {
		if j.Base == nil {
			j.Base = defaultBase
		} else {
			j.Base.ID = id
			j.Base.TenantID = tenantID
		}
		return j.copy()
	}
	return &Job{Base: defaultBase}
}

func (j *Job) copy() (n *Job) {
	n = NewJob("", "", nil)
	if j.Base != nil {
		n.Base = j.Base.copy()
	}
	n.ProtocolType = j.ProtocolType
	n.ProtocolConnectionID = j.ProtocolConnectionID
	n.ProtocolCredentialID = j.ProtocolCredentialID
	n.ProtocolProofID = j.ProtocolProofID
	n.ProtocolMessageID = j.ProtocolMessageID
	n.ConnectionID = j.ConnectionID
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

func (j *JobOutput) ToEdges() *model.JobOutput {
	var connection *model.PairwiseEdge
	if j.Connection != nil {
		connection = j.Connection.ToEdge()
	}
	var credential *model.CredentialEdge
	if j.Credential != nil {
		credential = j.Credential.ToEdge()
	}
	var proof *model.ProofEdge
	if j.Proof != nil {
		proof = j.Proof.ToEdge()
	}
	var message *model.BasicMessageEdge
	if j.Message != nil {
		message = j.Message.ToEdge()
	}
	return &model.JobOutput{
		Connection: connection,
		Credential: credential,
		Proof:      proof,
		Message:    message,
	}
}
