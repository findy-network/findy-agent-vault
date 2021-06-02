package model

import (
	"time"

	"github.com/findy-network/findy-agent-vault/graph/model"
	"github.com/findy-network/findy-agent-vault/paginator"
	"github.com/golang/glog"
)

type Proofs struct {
	Proofs          []*Proof
	HasNextPage     bool
	HasPreviousPage bool
}

type Proof struct {
	*base
	ConnectionID  string
	Role          model.ProofRole         `faker:"oneof: PROVER, PROVER"`
	Attributes    []*model.ProofAttribute `faker:"proofAttributes"`
	Values        []*model.ProofValue     `faker:"-"`
	InitiatedByUs bool
	Result        bool
	Provable      *time.Time `faker:"-"`
	Approved      *time.Time `faker:"-"`
	Verified      *time.Time `faker:"-"`
	Failed        *time.Time `faker:"-"`
	Archived      *time.Time `faker:"-"`
}

func NewProof(tenantID string, p *Proof) *Proof {
	defaultBase := &base{TenantID: tenantID}
	if p != nil {
		if p.base == nil {
			p.base = defaultBase
		} else {
			p.base.TenantID = tenantID
		}
		return p.copy()
	}
	return &Proof{base: defaultBase}
}

func (p *Proof) copy() (n *Proof) {
	n = NewProof("", nil)

	attributes := make([]*model.ProofAttribute, len(p.Attributes))
	for index := range p.Attributes {
		attributes[index] = &model.ProofAttribute{
			ID:        p.Attributes[index].ID,
			Name:      p.Attributes[index].Name,
			CredDefID: p.Attributes[index].CredDefID,
		}
	}

	values := make([]*model.ProofValue, len(p.Values))
	for index := range p.Values {
		values[index] = &model.ProofValue{
			ID:          p.Values[index].ID,
			AttributeID: p.Values[index].AttributeID,
			Value:       p.Values[index].Value,
		}
	}

	if p.base != nil {
		n.base = p.base.copy()
	}
	n.ConnectionID = p.ConnectionID
	n.Role = p.Role
	n.InitiatedByUs = p.InitiatedByUs
	n.Result = p.Result
	n.Provable = copyTime(p.Provable)
	n.Approved = copyTime(p.Approved)
	n.Verified = copyTime(p.Verified)
	n.Failed = copyTime(p.Failed)
	n.Attributes = attributes
	n.Values = values
	n.Archived = copyTime(p.Archived)

	return n
}

func (p *Proof) ToEdge() *model.ProofEdge {
	cursor := paginator.CreateCursor(p.Cursor, model.Proof{})
	return &model.ProofEdge{
		Cursor: cursor,
		Node:   p.ToNode(),
	}
}

func (p *Proof) ToNode() *model.Proof {
	return &model.Proof{
		ID:            p.ID,
		Role:          p.Role,
		Attributes:    p.Attributes,
		Values:        p.Values,
		InitiatedByUs: p.InitiatedByUs,
		Result:        p.Result,
		ApprovedMs:    timeToStringPtr(p.Approved),
		VerifiedMs:    timeToStringPtr(p.Verified),
		CreatedMs:     timeToString(&p.Created),
	}
}

func (p *Proof) Description() string {
	if p.Verified != nil {
		switch p.Role {
		case model.ProofRoleVerifier:
			return "Verified credential"
		case model.ProofRoleProver:
			return "Proved credential"
		}
	} else if p.Approved != nil {
		return "Approved proof"
	}
	switch p.Role {
	case model.ProofRoleVerifier:
		return "Received proof offer"
	case model.ProofRoleProver:
		return "Received proof request"
	}

	glog.Errorf("invalid role %s for proof", p.Role)
	return ""
}

func (p *Proofs) ToConnection(id *string) *model.ProofConnection {
	totalCount := len(p.Proofs)

	edges := make([]*model.ProofEdge, totalCount)
	nodes := make([]*model.Proof, totalCount)
	for index, connection := range p.Proofs {
		edge := connection.ToEdge()
		edges[index] = edge
		nodes[index] = edge.Node
	}

	var startCursor, endCursor *string
	if len(edges) > 0 {
		startCursor = &edges[0].Cursor
		endCursor = &edges[len(edges)-1].Cursor
	}
	return &model.ProofConnection{
		ConnectionID: id,
		Edges:        edges,
		Nodes:        nodes,
		PageInfo: &model.PageInfo{
			EndCursor:       endCursor,
			HasNextPage:     p.HasNextPage,
			HasPreviousPage: p.HasPreviousPage,
			StartCursor:     startCursor,
		},
	}
}
