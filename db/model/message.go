package model

import (
	"time"

	"github.com/findy-network/findy-agent-vault/graph/model"
	"github.com/findy-network/findy-agent-vault/paginator"
)

type Messages struct {
	Messages        []*Message
	HasNextPage     bool
	HasPreviousPage bool
}

type Message struct {
	Base
	ConnectionID string
	Message      string `faker:"sentence"`
	SentByMe     bool
	Delivered    *bool
	Archived     *time.Time `faker:"-"`
}

func (m *Message) ToEdge() *model.BasicMessageEdge {
	cursor := paginator.CreateCursor(m.Cursor, model.BasicMessage{})
	return &model.BasicMessageEdge{
		Cursor: cursor,
		Node:   m.ToNode(),
	}
}

func (m *Message) ToNode() *model.BasicMessage {
	var delivered *bool
	if m.Delivered != nil {
		d := *m.Delivered
		delivered = &d
	}
	return &model.BasicMessage{
		ID:        m.ID,
		Message:   m.Message,
		SentByMe:  m.SentByMe,
		Delivered: delivered,
		CreatedMs: timeToString(&m.Created),
	}
}

func (m *Message) Description() string {
	if m.SentByMe {
		return "Sent basic message"
	}
	return "Received basic message"
}

func (m *Messages) ToConnection(id *string) *model.BasicMessageConnection {
	totalCount := len(m.Messages)

	edges := make([]*model.BasicMessageEdge, totalCount)
	nodes := make([]*model.BasicMessage, totalCount)
	for index, connection := range m.Messages {
		edge := connection.ToEdge()
		edges[index] = edge
		nodes[index] = edge.Node
	}

	var startCursor, endCursor *string
	if len(edges) > 0 {
		startCursor = &edges[0].Cursor
		endCursor = &edges[len(edges)-1].Cursor
	}
	return &model.BasicMessageConnection{
		ConnectionID: id,
		Edges:        edges,
		Nodes:        nodes,
		PageInfo: &model.PageInfo{
			EndCursor:       endCursor,
			HasNextPage:     m.HasNextPage,
			HasPreviousPage: m.HasPreviousPage,
			StartCursor:     startCursor,
		},
	}
}
