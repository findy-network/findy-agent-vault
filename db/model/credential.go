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
	Base
	ConnectionID string
	Role         model.CredentialRole `faker:"oneof: HOLDER, HOLDER"`
	SchemaID     string
	CredDefID    string `faker:"oneof: credDefId1, credDefId2, credDefId3"`
	// TODO: can we avoid pointers with slices in gql interface?
	Attributes    []*model.CredentialValue `faker:"credentialAttributes"`
	InitiatedByUs bool
	Approved      time.Time `faker:"-"`
	Issued        time.Time `faker:"-"`
	Failed        time.Time `faker:"-"`
	Archived      time.Time `faker:"-"`
}

func (c *Credential) IsIssued() bool {
	return !c.Issued.IsZero()
}

func (c *Credential) ToEdge() *model.CredentialEdge {
	cursor := paginator.CreateCursor(c.Cursor, model.Credential{})
	return &model.CredentialEdge{
		Cursor: cursor,
		Node:   c.ToNode(),
	}
}

func (c *Credential) ToNode() *model.Credential {
	return &model.Credential{
		ID:            c.ID,
		Role:          c.Role,
		SchemaID:      c.SchemaID,
		CredDefID:     c.CredDefID,
		Attributes:    c.Attributes,
		InitiatedByUs: c.InitiatedByUs,
		ApprovedMs:    timeToStringPtr(&c.Approved),
		IssuedMs:      timeToStringPtr(&c.Issued),
		CreatedMs:     timeToString(&c.Created),
	}
}

func (c *Credential) Description() string {
	if !c.Issued.IsZero() {
		switch c.Role {
		case model.CredentialRoleIssuer:
			return "Issued credential"
		case model.CredentialRoleHolder:
			return "Received credential"
		}
	} else if !c.Approved.IsZero() {
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
