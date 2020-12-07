package model

import (
	"strconv"
	"time"

	"github.com/findy-network/findy-agent-vault/graph/model"
	graph "github.com/findy-network/findy-agent-vault/graph/model"
	"github.com/findy-network/findy-agent-vault/paginator"
)

type Connections struct {
	Connections     []*Connection
	HasNextPage     bool
	HasPreviousPage bool
}

type Connection struct {
	*base
	TenantID      string `faker:"tenantId"`
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

func (c *Connection) ToEdge() *graph.PairwiseEdge {
	cursor := paginator.CreateCursor(c.Cursor, graph.Pairwise{})
	return &graph.PairwiseEdge{
		Cursor: cursor,
		Node:   c.ToNode(),
	}
}

func (c *Connection) ToNode() *graph.Pairwise {
	approvedMs := ""
	if c.Approved != nil {
		approvedMs = strconv.FormatInt(c.Approved.UnixNano()/time.Millisecond.Nanoseconds(), 10)
	}
	return &graph.Pairwise{
		ID:            c.ID,
		OurDid:        c.OurDid,
		TheirDid:      c.TheirDid,
		TheirEndpoint: c.TheirEndpoint,
		TheirLabel:    c.TheirLabel,
		CreatedMs:     strconv.FormatInt(c.Created.UnixNano()/time.Millisecond.Nanoseconds(), 10),
		ApprovedMs:    approvedMs,
		Invited:       c.Invited,
	}
}

func ConnectionsToBatch(hasNextPage, hasPreviousPage bool, connections []*Connection) *graph.PairwiseConnection {
	totalCount := len(connections)

	edges := make([]*graph.PairwiseEdge, totalCount)
	nodes := make([]*graph.Pairwise, totalCount)
	for index, connection := range connections {
		edge := connection.ToEdge()
		edges[index] = edge
		nodes[index] = edge.Node
	}

	var startCursor, endCursor *string
	if len(edges) > 0 {
		startCursor = &edges[0].Cursor
		endCursor = &edges[len(edges)-1].Cursor
	}
	return &graph.PairwiseConnection{
		Edges: edges,
		Nodes: nodes,
		PageInfo: &model.PageInfo{
			EndCursor:       endCursor,
			HasNextPage:     hasNextPage,
			HasPreviousPage: hasPreviousPage,
			StartCursor:     startCursor,
		},
		TotalCount: totalCount, // TODO: total total count
	}
}
