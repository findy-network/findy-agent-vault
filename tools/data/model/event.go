package model

import (
	"strconv"

	"github.com/findy-network/findy-agent-vault/graph/model"
)

type InternalEvent struct {
	ID           string             `faker:"uuid_hyphenated"`
	Read         bool               `faker:"-"`
	Description  string             `faker:"sentence"`
	ProtocolType model.ProtocolType `faker:"oneof: model.ProtocolTypeNone, model.ProtocolTypeConnection, model.ProtocolTypeCredential, model.ProtocolTypeProof, model.ProtocolTypeBasicMessage"`
	Type         model.EventType    `faker:"oneof: model.EventTypeNotification, model.EventTypeAction"`
	PairwiseID   string             `faker:"eventPairwiseId"`
	CreatedMs    int64              `faker:"unix_time"`
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

func (e *InternalEvent) ToEdge(connections *Items) *model.EventEdge {
	cursor := CreateCursor(e.CreatedMs, model.Event{})
	return &model.EventEdge{
		Cursor: cursor,
		Node:   e.ToNode(connections),
	}
}

func (e *InternalEvent) ToNode(connections *Items) *model.Event {
	createdStr := strconv.FormatInt(e.CreatedMs, 10)
	return &model.Event{
		ID:          e.ID,
		Read:        e.Read,
		Description: e.Description,
		Protocol:    e.ProtocolType,
		Type:        e.Type,
		CreatedMs:   createdStr,
		Connection:  connections.PairwiseForID(e.PairwiseID),
	}
}

func (i *Items) EventForID(id string, connections *Items) *model.Event {
	var node *model.Event

	i.mutex.RLock()
	defer i.mutex.RUnlock()

	for _, item := range i.items {
		if item.Identifier() == id {
			event := item.Event()
			node = event.ToNode(connections)
			break
		}
	}

	return node
}

func (i *Items) EventConnection(after, before int, connections *Items) *model.EventConnection {
	i.mutex.RLock()
	result := i.items[after:before]
	totalCount := len(result)

	edges := make([]*model.EventEdge, totalCount)
	nodes := make([]*model.Event, totalCount)
	for index, event := range result {
		node := event.Event().ToNode(connections)
		edges[index] = &model.EventEdge{
			Cursor: CreateCursor(event.Event().CreatedMs, model.Event{}),
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
	c := &model.EventConnection{
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
