package model

import (
	"strconv"

	"github.com/findy-network/findy-agent-vault/graph/model"
)

type InternalMessage struct {
	ID         string `faker:"uuid_hyphenated"`
	Message    string `faker:"sentence"`
	PairwiseID string `faker:"pairwiseId"`
	SentByMe   bool
	Delivered  *bool
	CreatedMs  int64 `faker:"created"`
}

func (m *InternalMessage) Created() int64 {
	return m.CreatedMs
}

func (m *InternalMessage) Identifier() string {
	return m.ID
}

func (m *InternalMessage) Pairwise() *InternalPairwise {
	panic("Message is not connection")
}

func (m *InternalMessage) Event() *InternalEvent {
	panic("Message is not event")
}

func (m *InternalMessage) Job() *InternalJob {
	panic("Message is not job")
}

func (m *InternalMessage) BasicMessage() *InternalMessage {
	return m
}

func (m *InternalMessage) ToEdge() *model.BasicMessageEdge {
	cursor := CreateCursor(m.CreatedMs, model.BasicMessage{})
	return &model.BasicMessageEdge{
		Cursor: cursor,
		Node:   m.ToNode(),
	}
}

func (m *InternalMessage) ToNode() *model.BasicMessage {
	return &model.BasicMessage{
		ID:        m.ID,
		Message:   m.Message,
		SentByMe:  m.SentByMe,
		Delivered: m.Delivered,
		CreatedMs: strconv.FormatInt(m.CreatedMs, 10),
	}
}

func (i *Items) MessagePairwiseID(id string) (connectionID *string) {
	i.mutex.RLock()
	defer i.mutex.RUnlock()

	if id == "" {
		return
	}

	for _, item := range i.items {
		if item.Identifier() == id {
			c := item.BasicMessage().PairwiseID
			connectionID = &c
			break
		}
	}

	return
}

func (i *Items) MessageForID(id string) (edge *model.BasicMessageEdge) {
	i.mutex.RLock()
	defer i.mutex.RUnlock()

	if id == "" {
		return
	}

	for _, item := range i.items {
		if item.Identifier() == id {
			edge = item.BasicMessage().ToEdge()
			break
		}
	}

	return
}

func (i *Items) MessageConnection(after, before int) *model.BasicMessageConnection {
	i.mutex.RLock()
	result := i.items[after:before]
	totalCount := len(result)

	edges := make([]*model.BasicMessageEdge, totalCount)
	nodes := make([]*model.BasicMessage, totalCount)
	for index, m := range result {
		edge := m.BasicMessage().ToEdge()
		edges[index] = edge
		nodes[index] = edge.Node
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
	p := &model.BasicMessageConnection{
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
	return p
}
