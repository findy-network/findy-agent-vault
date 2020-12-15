package model

import (
	"time"

	"github.com/findy-network/findy-agent-vault/graph/model"
	"github.com/findy-network/findy-agent-vault/paginator"
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
	InitiatedByUs bool
	Result        bool
	Approved      *time.Time `faker:"-"`
	Verified      *time.Time `faker:"-"`
	Failed        *time.Time `faker:"-"`
}

func NewProof(p *Proof) *Proof {
	if p != nil {
		return p.copy()
	}
	return &Proof{base: &base{}}
}

func (p *Proof) copy() (n *Proof) {
	n = NewProof(nil)

	attributes := make([]*model.ProofAttribute, len(p.Attributes))
	for index := range p.Attributes {
		attributes[index] = &model.ProofAttribute{
			ID:        p.Attributes[index].ID,
			Name:      p.Attributes[index].Name,
			Value:     p.Attributes[index].Value,
			CredDefID: p.Attributes[index].CredDefID,
		}
	}

	if p.base != nil {
		n.base = p.base.copy()
	}
	n.ConnectionID = p.ConnectionID
	n.Role = p.Role
	n.InitiatedByUs = p.InitiatedByUs
	n.Result = p.Result
	n.Approved = copyTime(p.Approved)
	n.Verified = copyTime(p.Verified)
	n.Failed = copyTime(p.Failed)
	n.Attributes = attributes

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
	approvedMs := ""
	verifiedMs := ""
	if p.Approved != nil {
		approvedMs = timeToString(p.Approved)
	}
	if p.Verified != nil {
		verifiedMs = timeToString(p.Verified)
	}
	return &model.Proof{
		ID:            p.ID,
		Role:          p.Role,
		Attributes:    p.Attributes,
		InitiatedByUs: p.InitiatedByUs,
		Result:        p.Result,
		ApprovedMs:    &approvedMs,
		VerifiedMs:    &verifiedMs,
		CreatedMs:     timeToString(&p.Created),
	}
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
