package model

import (
	"time"

	"github.com/findy-network/findy-agent-vault/graph/model"
	"github.com/findy-network/findy-agent-vault/paginator"
)

type Credentials struct {
	Credentials     []*Credential
	HasNextPage     bool
	HasPreviousPage bool
}

type Credential struct {
	*base
	TenantID      string `faker:"tenantId"`
	ConnectionID  string
	Role          model.CredentialRole `faker:"oneof: HOLDER, HOLDER"`
	SchemaID      string
	CredDefID     string
	Attributes    []*model.CredentialValue `faker:"credentialAttributes"`
	InitiatedByUs bool
	Approved      *time.Time `faker:"-"`
	Issued        *time.Time `faker:"-"`
	Failed        *time.Time `faker:"-"`
	Cursor        uint64
}

func NewCredential() *Credential { return &Credential{base: &base{}} }

// Note: not deep copy!
func (c *Credential) Copy() (n *Credential) {
	n = NewCredential()
	n.TenantID = c.TenantID
	n.ConnectionID = c.ConnectionID
	n.Role = c.Role
	n.SchemaID = c.SchemaID
	n.CredDefID = c.CredDefID
	n.InitiatedByUs = c.InitiatedByUs
	n.Approved = c.Approved
	n.Issued = c.Issued
	n.Failed = c.Failed
	n.Attributes = c.Attributes

	return n
}

func (c *Credential) ToEdge() *model.CredentialEdge {
	cursor := paginator.CreateCursor(c.Cursor, model.Credential{})
	return &model.CredentialEdge{
		Cursor: cursor,
		Node:   c.ToNode(),
	}
}

func (c *Credential) ToNode() *model.Credential {
	approvedMs := ""
	issuedMs := ""
	if c.Approved != nil {
		approvedMs = timeToString(c.Approved)
	}
	if c.Issued != nil {
		issuedMs = timeToString(c.Issued)
	}
	return &model.Credential{
		ID:            c.ID,
		Role:          c.Role,
		SchemaID:      c.SchemaID,
		CredDefID:     c.CredDefID,
		Attributes:    c.Attributes,
		InitiatedByUs: c.InitiatedByUs,
		ApprovedMs:    &approvedMs,
		IssuedMs:      &issuedMs,
		CreatedMs:     timeToString(&c.Created),
	}
}

func (c *Credentials) ToConnection() *model.CredentialConnection {
	totalCount := len(c.Credentials)

	edges := make([]*model.CredentialEdge, totalCount)
	nodes := make([]*model.Credential, totalCount)
	for index, connection := range c.Credentials {
		edge := connection.ToEdge()
		edges[index] = edge
		nodes[index] = edge.Node
	}

	var startCursor, endCursor *string
	if len(edges) > 0 {
		startCursor = &edges[0].Cursor
		endCursor = &edges[len(edges)-1].Cursor
	}
	return &model.CredentialConnection{
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
