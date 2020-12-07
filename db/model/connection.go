package model

import (
	"strconv"
	"time"

	graph "github.com/findy-network/findy-agent-vault/graph/model"
)

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
	cursor := CreateCursor(c.Cursor, graph.Pairwise{})
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
