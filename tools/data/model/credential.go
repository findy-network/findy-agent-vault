package model

import (
	"strconv"

	"github.com/findy-network/findy-agent-vault/graph/model"
)

type CredentialItems struct {
	*Items
}

type InternalCredential struct {
	*BaseObject
	Role          model.CredentialRole
	SchemaID      string
	CredDefID     string
	Attributes    []*model.CredentialValue
	InitiatedByUs bool
	ApprovedMs    *int64
	IssuedMs      *int64
	FailedMs      *int64
	PairwiseID    string `faker:"pairwiseId"`
}

func (c *InternalCredential) Description() string {
	if c.IssuedMs != nil {
		switch c.Role {
		case model.CredentialRoleIssuer:
			return "Issued credential"
		case model.CredentialRoleHolder:
			return "Received credential"
		}
	} else if c.ApprovedMs != nil {
		return "Approved credential"
	}
	switch c.Role {
	case model.CredentialRoleIssuer:
		return "Received credential request"
	case model.CredentialRoleHolder:
		return "Received credential offer"
	}
	return ""
}

func (c *InternalCredential) Status() *ProtocolStatus {
	status := model.JobStatusWaiting
	result := model.JobResultNone
	if c.FailedMs != nil {
		status = model.JobStatusComplete
		result = model.JobResultFailure
	} else if c.ApprovedMs == nil && c.IssuedMs == nil {
		status = model.JobStatusPending
	} else if c.IssuedMs != nil {
		status = model.JobStatusComplete
		result = model.JobResultSuccess
	}

	return &ProtocolStatus{
		Status:      status,
		Result:      result,
		Description: c.Description(),
	}
}

func (c *InternalCredential) Credential() *InternalCredential {
	return c
}

func (c *InternalCredential) Copy() *InternalCredential {
	var approvedMs, issuedMs *int64
	if c.ApprovedMs != nil {
		a := *c.ApprovedMs
		approvedMs = &a
	}
	if c.IssuedMs != nil {
		i := *c.IssuedMs
		issuedMs = &i
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
		Role:          c.Role,
		SchemaID:      c.SchemaID,
		CredDefID:     c.CredDefID,
		Attributes:    values,
		InitiatedByUs: c.InitiatedByUs,
		ApprovedMs:    approvedMs,
		IssuedMs:      issuedMs,
		PairwiseID:    c.PairwiseID,
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
	var approvedMs, issuedMs *string
	if cred.ApprovedMs != nil {
		a := strconv.FormatInt(*cred.ApprovedMs, 10)
		approvedMs = &a
	}
	if cred.IssuedMs != nil {
		i := strconv.FormatInt(*cred.IssuedMs, 10)
		issuedMs = &i
	}
	return &model.Credential{
		ID:            cred.ID,
		Role:          cred.Role,
		SchemaID:      cred.SchemaID,
		CredDefID:     cred.CredDefID,
		Attributes:    cred.Attributes,
		InitiatedByUs: cred.InitiatedByUs,
		ApprovedMs:    approvedMs,
		IssuedMs:      issuedMs,
		CreatedMs:     strconv.FormatInt(cred.CreatedMs, 10),
	}
}

func (i *CredentialItems) Objects() *Items {
	return i.Items
}

func (i *CredentialItems) CredentialPairwiseID(id string) (connectionID *string) {
	i.mutex.RLock()
	defer i.mutex.RUnlock()

	if id == "" {
		return
	}

	for _, item := range i.items {
		if item.Identifier() == id {
			c := item.Credential().PairwiseID
			connectionID = &c
			break
		}
	}

	return
}

func (i *CredentialItems) CredentialForID(id string) (edge *model.CredentialEdge) {
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

func (i *CredentialItems) CredentialConnection(after, before int) *model.CredentialConnection {
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

func (i *CredentialItems) UpdateCredential(id string, approvedMs, issuedMs, failedMs *int64) *ProtocolStatus {
	i.mutex.Lock()
	defer i.mutex.Unlock()

	for _, item := range i.items {
		if item.Identifier() != id {
			continue
		}
		cred := item.Credential()
		if approvedMs != nil {
			cred.ApprovedMs = approvedMs
		}
		if issuedMs != nil {
			cred.IssuedMs = issuedMs
		}
		if failedMs != nil {
			cred.FailedMs = failedMs
		}
		return cred.Status()
	}

	return nil
}
