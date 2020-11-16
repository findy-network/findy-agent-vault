package model

import (
	"strconv"

	"github.com/findy-network/findy-agent-vault/graph/model"
)

type CredentialItems interface {
	CredentialConnection(after, before int) *model.CredentialConnection
	CredentialForID(id string) *model.CredentialEdge
	Objects() *Items
}

func (i *Items) Credentials() CredentialItems { return &credentialItems{i} }

type credentialItems struct{ *Items }

type InternalCredential struct {
	*BaseObject
	CredentialRole model.CredentialRole
	SchemaID       string
	CredDefID      string
	Attributes     []*model.CredentialValue
	InitiatedByUs  bool
	ApprovedMs     *int64
	PairwiseID     string `faker:"pairwiseId"`
}

func (c *InternalCredential) Credential() *InternalCredential {
	return c
}

func (c *InternalCredential) Copy() *InternalCredential {
	var approvedMs *int64
	if c.ApprovedMs != nil {
		a := *c.ApprovedMs
		approvedMs = &a
	}
	values := make([]*model.CredentialValue, 0)
	for i := range c.Attributes {
		v := *c.Attributes[i]
		values = append(values, &v)
	}
	newCred := &InternalCredential{
		BaseObject: &BaseObject{
			ID:        c.ID,
			CreatedMs: c.CreatedMs,
		},
		CredentialRole: c.CredentialRole,
		SchemaID:       c.SchemaID,
		CredDefID:      c.CredDefID,
		Attributes:     values,
		InitiatedByUs:  c.InitiatedByUs,
		ApprovedMs:     approvedMs,
		PairwiseID:     c.PairwiseID,
	}
	return newCred
}

func (c *InternalCredential) ToEdge() *model.CredentialEdge {
	cursor := CreateCursor(c.CreatedMs, model.Credential{})
	return &model.CredentialEdge{
		Cursor: cursor,
		Node:   c.ToNode(),
	}
}

func (c *InternalCredential) ToNode() *model.Credential {

	cred := c.Copy()
	var approvedMs *string
	if cred.ApprovedMs != nil {
		a := strconv.FormatInt(*cred.ApprovedMs, 10)
		approvedMs = &a
	}
	return &model.Credential{
		ID:            cred.ID,
		Role:          cred.CredentialRole,
		SchemaID:      cred.SchemaID,
		CredDefID:     cred.CredDefID,
		Attributes:    cred.Attributes,
		InitiatedByUs: cred.InitiatedByUs,
		ApprovedMs:    approvedMs,
		CreatedMs:     strconv.FormatInt(cred.CreatedMs, 10),
	}
}

func (i *credentialItems) Objects() *Items {
	return i.Items
}

func (i *credentialItems) CredentialForID(id string) (edge *model.CredentialEdge) {
	i.mutex.RLock()
	defer i.mutex.RUnlock()

	if id == "" {
		return
	}

	for _, item := range i.items {
		if item.Identifier() == id {
			edge = item.Credential().ToEdge()
			break
		}
	}

	return
}

func (i *credentialItems) CredentialConnection(after, before int) *model.CredentialConnection {
	i.mutex.RLock()
	result := i.items[after:before]
	totalCount := len(result)

	edges := make([]*model.CredentialEdge, totalCount)
	nodes := make([]*model.Credential, totalCount)
	for index, pairwise := range result {
		edge := pairwise.Credential().ToEdge()
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
	p := &model.CredentialConnection{
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
