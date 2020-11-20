package model

import (
	"fmt"
	"strconv"

	"github.com/findy-network/findy-agent-vault/graph/model"
)

type ProofItems struct {
	*Items
}

type InternalProof struct {
	*BaseObject
	Role          model.ProofRole `faker:"oneof: PROVER, PROVER"`
	Attributes    []*model.ProofAttribute
	InitiatedByUs bool
	Result        bool
	VerifiedMs    *int64
	ApprovedMs    *int64
	FailedMs      *int64
	PairwiseID    string `faker:"pairwiseId"`
}

func (p *InternalProof) Description() string {
	if p.VerifiedMs != nil {
		switch p.Role {
		case model.ProofRoleVerifier:
			return "Verified credential"
		case model.ProofRoleProver:
			return "Proved credential"
		}
	} else if p.ApprovedMs != nil {
		return "Approved proof"
	}
	switch p.Role {
	case model.ProofRoleVerifier:
		return "Received proof offer"
	case model.ProofRoleProver:
		return "Received proof request"
	}
	panic(fmt.Errorf("invalid role %s for proof", p.Role))
}

func (p *InternalProof) Status() *ProtocolStatus {
	status := model.JobStatusWaiting
	result := model.JobResultNone
	if p.FailedMs != nil {
		status = model.JobStatusComplete
		result = model.JobResultFailure
	} else if p.ApprovedMs == nil && p.VerifiedMs == nil {
		status = model.JobStatusPending
	} else if p.VerifiedMs != nil {
		status = model.JobStatusComplete
		result = model.JobResultSuccess
	}

	return &ProtocolStatus{
		Status:      status,
		Result:      result,
		Description: p.Description(),
	}
}

func (p *InternalProof) Proof() *InternalProof {
	return p
}

func (p *InternalProof) Copy() *InternalProof {
	var approvedMs, verifiedMs *int64
	if p.ApprovedMs != nil {
		a := *p.ApprovedMs
		approvedMs = &a
	}
	if p.VerifiedMs != nil {
		v := *p.VerifiedMs
		verifiedMs = &v
	}
	attributes := make([]*model.ProofAttribute, 0)
	for i := range p.Attributes {
		a := *p.Attributes[i]
		attributes = append(attributes, &a)
	}
	newProof := &InternalProof{
		BaseObject: &BaseObject{
			ID:        p.ID,
			CreatedMs: p.CreatedMs,
		},
		Role:          p.Role,
		Attributes:    attributes,
		InitiatedByUs: p.InitiatedByUs,
		Result:        p.Result,
		VerifiedMs:    verifiedMs,
		ApprovedMs:    approvedMs,
		PairwiseID:    p.PairwiseID,
	}
	return newProof
}

func (p *InternalProof) ToEdge() *model.ProofEdge {
	cursor := CreateCursor(p.CreatedMs, model.Proof{})
	return &model.ProofEdge{
		Cursor: cursor,
		Node:   p.ToNode(),
	}
}

func (p *InternalProof) ToNode() *model.Proof {
	var approvedMs, verifiedMs *string
	if p.ApprovedMs != nil {
		a := strconv.FormatInt(*p.ApprovedMs, 10)
		approvedMs = &a
	}
	if p.VerifiedMs != nil {
		v := strconv.FormatInt(*p.VerifiedMs, 10)
		verifiedMs = &v
	}
	return &model.Proof{
		ID:            p.ID,
		Role:          p.Role,
		Attributes:    p.Attributes,
		InitiatedByUs: p.InitiatedByUs,
		Result:        p.Result,
		VerifiedMs:    verifiedMs,
		ApprovedMs:    approvedMs,
		CreatedMs:     strconv.FormatInt(p.CreatedMs, 10),
	}
}

func (i *ProofItems) Objects() *Items {
	return i.Items
}

func (i *ProofItems) ProofPairwiseID(id string) (connectionID *string) {
	i.mutex.RLock()
	defer i.mutex.RUnlock()

	if id == "" {
		return
	}

	for _, item := range i.items {
		if item.Identifier() == id {
			c := item.Proof().PairwiseID
			connectionID = &c
			break
		}
	}

	return
}

func (i *ProofItems) ProofForID(id string) (edge *model.ProofEdge) {
	i.mutex.RLock()
	defer i.mutex.RUnlock()

	if id == "" {
		return
	}

	for _, item := range i.items {
		if item.Identifier() == id {
			edge = item.Proof().ToEdge()
			break
		}
	}

	return
}

func (i *ProofItems) ProofConnection(after, before int) *model.ProofConnection {
	i.mutex.RLock()
	result := i.items[after:before]
	totalCount := len(result)

	edges := make([]*model.ProofEdge, totalCount)
	nodes := make([]*model.Proof, totalCount)
	for index, item := range result {
		edge := item.Proof().ToEdge()
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
	p := &model.ProofConnection{
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

func (i *ProofItems) UpdateProof(id string, result *bool, verifiedMs, approvedMs, failedMs *int64) *ProtocolStatus {
	i.mutex.Lock()
	defer i.mutex.Unlock()

	for _, item := range i.items {
		if item.Identifier() != id {
			continue
		}
		proof := item.Proof()
		if result != nil {
			proof.Result = *result
		}
		if verifiedMs != nil {
			proof.VerifiedMs = verifiedMs
		}
		if approvedMs != nil {
			proof.ApprovedMs = approvedMs
		}
		if failedMs != nil {
			proof.FailedMs = failedMs
		}
		return proof.Status()
	}

	return nil
}
