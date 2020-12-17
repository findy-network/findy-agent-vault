package mock

import (
	"errors"
	"time"

	"github.com/bxcodec/faker/v3"
	"github.com/findy-network/findy-agent-vault/db/model"
	"github.com/findy-network/findy-agent-vault/paginator"
)

type mockProof struct {
	*base
	proof *model.Proof
}

func (p *mockProof) Created() uint64 {
	return model.TimeToCursor(&p.proof.Created)
}

func (p *mockProof) Identifier() string {
	return p.proof.ID
}

func newProof(p *model.Proof) *mockProof {
	var proof *model.Proof
	if p != nil {
		proof = model.NewProof(p.TenantID, p)
	}
	return &mockProof{base: &base{}, proof: proof}
}

func (p *mockProof) Copy() apiObject {
	return newProof(p.proof)
}

func (p *mockProof) Proof() *model.Proof {
	return p.proof
}

func (m *mockData) AddProof(p *model.Proof) (*model.Proof, error) {
	agent := m.agents.get(p.TenantID)

	n := model.NewProof(p.TenantID, p)
	n.ID = faker.UUIDHyphenated()
	n.Created = time.Now().UTC()
	n.Cursor = model.TimeToCursor(&n.Created)
	for index := range n.Attributes {
		n.Attributes[index].ID = faker.UUIDHyphenated()
	}
	agent.proofs.append(newProof(n))
	return n, nil
}

func (m *mockData) UpdateProof(p *model.Proof) (*model.Proof, error) {
	agent := m.agents.get(p.TenantID)

	object := agent.proofs.objectForID(p.ID)
	if object == nil {
		return nil, errors.New("not found proof for id: " + p.ID)
	}
	updated := object.Copy()
	proof := updated.Proof()
	proof.Approved = p.Approved
	proof.Verified = p.Verified
	proof.Failed = p.Failed

	if !agent.proofs.replaceObjectForID(p.ID, updated) {
		return nil, errors.New("not found proof for id: " + p.ID)
	}
	return updated.Proof(), nil
}

func (m *mockData) GetProof(id, tenantID string) (*model.Proof, error) {
	agent := m.agents.get(tenantID)

	p := agent.proofs.objectForID(id)
	if p == nil {
		return nil, errors.New("not found proof for id: " + id)
	}
	return p.Proof(), nil
}

func filterProof(item apiObject) bool {
	p := item.Proof()
	return p.Verified != nil
}

func proofConnectionFilter(id string) func(item apiObject) bool {
	return func(item apiObject) bool {
		p := item.Proof()
		if p.Verified != nil && p.ConnectionID == id {
			return true
		}
		return false
	}
}

func (m *mockItems) getProofs(
	info *paginator.BatchInfo,
	filter func(item apiObject) bool,
) (connections *model.Proofs, err error) {
	state, hasNextPage, hasPreviousPage := m.proofs.getObjects(info, filter)
	res := make([]*model.Proof, len(state.objects))
	for i := range state.objects {
		res[i] = state.objects[i].Copy().Proof()
	}

	c := &model.Proofs{
		Proofs:          res,
		HasNextPage:     hasNextPage,
		HasPreviousPage: hasPreviousPage,
	}
	return c, nil
}

func (m *mockData) GetProofs(
	info *paginator.BatchInfo,
	tenantID string,
	connectionID *string,
) (connections *model.Proofs, err error) {
	agent := m.agents.get(tenantID)

	if connectionID == nil {
		return agent.getProofs(info, filterProof)
	}
	return agent.getProofs(info, proofConnectionFilter(*connectionID))
}

func (m *mockData) GetProofCount(tenantID string, connectionID *string) (int, error) {
	agent := m.agents.get(tenantID)

	if connectionID == nil {
		return agent.proofs.count(filterProof), nil
	}
	return agent.proofs.count(proofConnectionFilter(*connectionID)), nil
}

func (m *mockData) GetConnectionForProof(id, tenantID string) (*model.Connection, error) {
	proof, err := m.GetProof(id, tenantID)
	if err != nil {
		return nil, err
	}
	return m.GetConnection(proof.ConnectionID, tenantID)
}
