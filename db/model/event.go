package model

import (
	"github.com/findy-network/findy-agent-vault/graph/model"
	"github.com/findy-network/findy-agent-vault/paginator"
)

type Events struct {
	Events          []*Event
	HasNextPage     bool
	HasPreviousPage bool
}

type Event struct {
	*base
	Read         bool    `faker:"-"`
	Description  string  `faker:"sentence"`
	JobID        *string `faker:"-"`
	ConnectionID *string `faker:"-"`
}

func NewEvent(e *Event) *Event {
	if e != nil {
		return e.copy()
	}
	return &Event{base: &base{}}
}

func (e *Event) copy() (n *Event) {
	n = NewEvent(nil)
	if e.base != nil {
		n.base = e.base.Copy()
	}
	n.Read = e.Read
	n.Description = e.Description
	if e.JobID != nil {
		jobID := *e.JobID
		n.JobID = &jobID
	}
	if e.ConnectionID != nil {
		connectionID := *e.ConnectionID
		n.ConnectionID = &connectionID
	}
	return n
}

func (e *Event) ToEdge() *model.EventEdge {
	cursor := paginator.CreateCursor(e.Cursor, model.Event{})
	return &model.EventEdge{
		Cursor: cursor,
		Node:   e.ToNode(),
	}
}

func (e *Event) ToNode() *model.Event {
	return &model.Event{
		ID:          e.ID,
		Read:        e.Read,
		Description: e.Description,
		CreatedMs:   timeToString(&e.Created),
	}
}

func (e *Events) ToConnection() *model.EventConnection {
	totalCount := len(e.Events)

	edges := make([]*model.EventEdge, totalCount)
	nodes := make([]*model.Event, totalCount)
	for index, event := range e.Events {
		edge := event.ToEdge()
		edges[index] = edge
		nodes[index] = edge.Node
	}

	var startCursor, endCursor *string
	if len(edges) > 0 {
		startCursor = &edges[0].Cursor
		endCursor = &edges[len(edges)-1].Cursor
	}
	return &model.EventConnection{
		Edges: edges,
		Nodes: nodes,
		PageInfo: &model.PageInfo{
			EndCursor:       endCursor,
			HasNextPage:     e.HasNextPage,
			HasPreviousPage: e.HasPreviousPage,
			StartCursor:     startCursor,
		},
	}
}
