package model

import (
	"strconv"

	"github.com/findy-network/findy-agent-vault/graph/model"
)

type ProofItems interface {
	ProofConnection(after, before int) *model.ProofConnection
	ProofForID(id string) *model.ProofEdge
	ProofPairwiseID(id string) *string
	UpdateProof(id string, result *bool, verifiedMs, approvedMs, failedMs *int64) *ProtocolStatus
	Objects() *Items
}

func (i *Items) Proofs() ProofItems { return &proofItems{i} }

type proofItems struct{ *Items }

type InternalProof struct {
	*BaseObject
	Role          model.ProofRole
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
	return ""
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
	newCred := &InternalProof{
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
	return newCred
}

func (p *InternalProof) ToEdge() *model.ProofEdge {
	cursor := CreateCursor(p.CreatedMs, model.Proof{})
	return &model.ProofEdge{
		Cursor: cursor,
		Node:   p.ToNode(),
	}
}

func (p *InternalProof) ToNode() *model.Proof {
	proof := p.Copy()
	var approvedMs, verifiedMs *string
	if proof.ApprovedMs != nil {
		a := strconv.FormatInt(*proof.ApprovedMs, 10)
		approvedMs = &a
	}
	if proof.VerifiedMs != nil {
		v := strconv.FormatInt(*proof.VerifiedMs, 10)
		verifiedMs = &v
	}
	return &model.Proof{
		ID:            proof.ID,
		Role:          proof.Role,
		Attributes:    proof.Attributes,
		InitiatedByUs: proof.InitiatedByUs,
		Result:        p.Result,
		VerifiedMs:    verifiedMs,
		ApprovedMs:    approvedMs,
		CreatedMs:     strconv.FormatInt(proof.CreatedMs, 10),
	}
}

func (i *proofItems) Objects() *Items {
	return i.Items
}

func (i *proofItems) ProofPairwiseID(id string) (connectionID *string) {
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

func (i *proofItems) ProofForID(id string) (edge *model.ProofEdge) {
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

func (i *proofItems) ProofConnection(after, before int) *model.ProofConnection {
	i.mutex.RLock()
	result := i.items[after:before]
	totalCount := len(result)

	edges := make([]*model.ProofEdge, totalCount)
	nodes := make([]*model.Proof, totalCount)
	for index, pairwise := range result {
		edge := pairwise.Proof().ToEdge()
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

func (i *proofItems) UpdateProof(id string, result *bool, verifiedMs, approvedMs, failedMs *int64) *ProtocolStatus {
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
