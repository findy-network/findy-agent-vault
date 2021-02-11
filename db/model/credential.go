package model

import (
	"time"

	"github.com/findy-network/findy-agent-vault/graph/model"
	"github.com/findy-network/findy-agent-vault/paginator"
	"github.com/golang/glog"
)

type Credentials struct {
	Credentials     []*Credential
	HasNextPage     bool
	HasPreviousPage bool
}

type Credential struct {
	*base
	ConnectionID  string
	Role          model.CredentialRole `faker:"oneof: HOLDER, HOLDER"`
	SchemaID      string
	CredDefID     string
	Attributes    []*model.CredentialValue `faker:"credentialAttributes"`
	InitiatedByUs bool
	Approved      *time.Time `faker:"-"`
	Issued        *time.Time `faker:"-"`
	Failed        *time.Time `faker:"-"`
	Archived      *time.Time `faker:"-"`
}

func NewCredential(tenantID string, c *Credential) *Credential {
	defaultBase := &base{TenantID: tenantID}
	if c != nil {
		if c.base == nil {
			c.base = defaultBase
		} else {
			c.base.TenantID = tenantID
		}
		return c.copy()
	}
	return &Credential{base: defaultBase}
}

func (c *Credential) copy() (n *Credential) {
	n = NewCredential("", nil)

	attributes := make([]*model.CredentialValue, len(c.Attributes))
	for index := range c.Attributes {
		attributes[index] = &model.CredentialValue{
			ID:    c.Attributes[index].ID,
			Name:  c.Attributes[index].Name,
			Value: c.Attributes[index].Value,
		}
	}

	if c.base != nil {
		n.base = c.base.copy()
	}
	n.ConnectionID = c.ConnectionID
	n.Role = c.Role
	n.SchemaID = c.SchemaID
	n.CredDefID = c.CredDefID
	n.InitiatedByUs = c.InitiatedByUs
	n.Approved = copyTime(c.Approved)
	n.Issued = copyTime(c.Issued)
	n.Failed = copyTime(c.Failed)
	n.Archived = copyTime(c.Archived)
	n.Attributes = attributes

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

func (c *Credential) Description() string {
	if c.Issued != nil {
		switch c.Role {
		case model.CredentialRoleIssuer:
			return "Issued credential"
		case model.CredentialRoleHolder:
			return "Received credential"
		}
	} else if c.Approved != nil {
		return "Approved credential"
	}

	switch c.Role {
	case model.CredentialRoleIssuer:
		return "Received credential request"
	case model.CredentialRoleHolder:
		return "Received credential offer"
	}

	glog.Errorf("invalid role %s for credential", c.Role)
	return ""
}

func (c *Credentials) ToConnection(id *string) *model.CredentialConnection {
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
		ConnectionID: id,
		Edges:        edges,
		Nodes:        nodes,
		PageInfo: &model.PageInfo{
			EndCursor:       endCursor,
			HasNextPage:     c.HasNextPage,
			HasPreviousPage: c.HasPreviousPage,
			StartCursor:     startCursor,
		},
	}
}
