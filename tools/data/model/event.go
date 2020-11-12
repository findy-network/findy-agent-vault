package model

import (
	"strconv"

	"github.com/findy-network/findy-agent-vault/graph/model"
)

type InternalEvent struct {
	ID          string  `faker:"uuid_hyphenated"`
	Read        bool    `faker:"-"`
	Description string  `faker:"sentence"`
	JobID       *string `faker:"-"`
	PairwiseID  *string `faker:"eventPairwiseId"`
	CreatedMs   int64   `faker:"unix_time"`
}

func (e *InternalEvent) Created() int64 {
	return e.CreatedMs
}

func (e *InternalEvent) Identifier() string {
	return e.ID
}

func (e *InternalEvent) Pairwise() *InternalPairwise {
	panic("Event is not pairwise")
}

func (e *InternalEvent) Event() *InternalEvent {
	return e
}

func (e *InternalEvent) Job() *InternalJob {
	panic("Event is not job")
}

func (e *InternalEvent) ToEdge(connections, jobs *Items) *model.EventEdge {
	cursor := CreateCursor(e.CreatedMs, model.Event{})
	return &model.EventEdge{
		Cursor: cursor,
		Node:   e.ToNode(connections, jobs),
	}
}

func (e *InternalEvent) ToNode(connections, jobs *Items) *model.Event {
	createdStr := strconv.FormatInt(e.CreatedMs, 10)
	var pw *model.Pairwise
	var job *model.JobEdge
	if e.PairwiseID != nil {
		if edge := connections.PairwiseForID(*e.PairwiseID); edge != nil {
			pw = edge.Node
		}
	}
	if e.JobID != nil {
		job = jobs.JobForID(*e.JobID, connections)
	}
	return &model.Event{
		ID:          e.ID,
		Read:        e.Read,
		Description: e.Description,
		CreatedMs:   createdStr,
		Connection:  pw,
		Job:         job,
	}
}

func (i *Items) EventForID(id string, connections, jobs *Items) (edge *model.EventEdge) {
	i.mutex.RLock()
	defer i.mutex.RUnlock()

	for _, item := range i.items {
		if item.Identifier() == id {
			event := item.Event()
			edge = event.ToEdge(connections, jobs)
			break
		}
	}

	return edge
}

func (i *Items) EventConnection(after, before int, connections, jobs *Items) *model.EventConnection {
	i.mutex.RLock()
	result := i.items[after:before]
	totalCount := len(result)

	edges := make([]*model.EventEdge, totalCount)
	nodes := make([]*model.Event, totalCount)
	for index, event := range result {
		node := event.Event().ToNode(connections, jobs)
		edges[index] = &model.EventEdge{
			Cursor: CreateCursor(event.Event().CreatedMs, model.Event{}),
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
	c := &model.EventConnection{
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

func (i *Items) MarkEventRead(id string) bool {
	i.mutex.Lock()
	defer i.mutex.Unlock()

	for _, item := range i.items {
		if item.Identifier() == id {
			event := item.Event()
			event.Read = true
			return true
		}
	}

	return false
}
