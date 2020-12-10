package model

import (
	"time"

	"github.com/findy-network/findy-agent-vault/graph/model"
	"github.com/findy-network/findy-agent-vault/paginator"
)

type Connections struct {
	Connections     []*Connection
	HasNextPage     bool
	HasPreviousPage bool
}

type Connection struct {
	*base
	TenantID      string
	OurDid        string
	TheirDid      string
	TheirEndpoint string `faker:"url"`
	TheirLabel    string `faker:"organisationLabel"`
	Invited       bool
	Approved      *time.Time
	Cursor        uint64
}

func NewConnection() *Connection { return &Connection{base: &base{}} }

func (c *Connection) Copy() (n *Connection) {
	n = NewConnection()
	n.TenantID = c.TenantID
	n.OurDid = c.OurDid
	n.TheirDid = c.TheirDid
	n.TheirEndpoint = c.TheirEndpoint
	n.TheirLabel = c.TheirLabel
	n.Invited = c.Invited
	return n
}

func (c *Connection) ToEdge() *model.PairwiseEdge {
	cursor := paginator.CreateCursor(c.Cursor, model.Pairwise{})
	return &model.PairwiseEdge{
		Cursor: cursor,
		Node:   c.ToNode(),
	}
}

func (c *Connection) ToNode() *model.Pairwise {
	approvedMs := ""
	if c.Approved != nil {
		approvedMs = timeToString(c.Approved)
	}
	return &model.Pairwise{
		ID:            c.ID,
		OurDid:        c.OurDid,
		TheirDid:      c.TheirDid,
		TheirEndpoint: c.TheirEndpoint,
		TheirLabel:    c.TheirLabel,
		CreatedMs:     timeToString(&c.Created),
		ApprovedMs:    approvedMs,
		Invited:       c.Invited,
	}
}

func (c *Connections) ToConnection() *model.PairwiseConnection {
	totalCount := len(c.Connections)

	edges := make([]*model.PairwiseEdge, totalCount)
	nodes := make([]*model.Pairwise, totalCount)
	for index, connection := range c.Connections {
		edge := connection.ToEdge()
		edges[index] = edge
		nodes[index] = edge.Node
	}

	var startCursor, endCursor *string
	if len(edges) > 0 {
		startCursor = &edges[0].Cursor
		endCursor = &edges[len(edges)-1].Cursor
	}
	return &model.PairwiseConnection{
		Edges: edges,
		Nodes: nodes,
		PageInfo: &model.PageInfo{
			EndCursor:       endCursor,
			HasNextPage:     c.HasNextPage,
			HasPreviousPage: c.HasPreviousPage,
			StartCursor:     startCursor,
		},
		TotalCount: totalCount, // TODO: total total count
	}
}
