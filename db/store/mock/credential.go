package mock

import (
	"errors"
	"time"

	"github.com/bxcodec/faker/v3"
	"github.com/findy-network/findy-agent-vault/db/model"
	"github.com/findy-network/findy-agent-vault/paginator"
)

type mockCredential struct {
	*base
	credential *model.Credential
}

func (c *mockCredential) Created() uint64 {
	return model.TimeToCursor(&c.credential.Created)
}

func (c *mockCredential) Identifier() string {
	return c.credential.ID
}

func newCredential(c *model.Credential) *mockCredential {
	var credential *model.Credential
	if c != nil {
		credential = model.NewCredential(c.TenantID, c)
	}
	return &mockCredential{base: &base{}, credential: credential}
}

func (c *mockCredential) Copy() apiObject {
	return newCredential(c.credential)
}

func (c *mockCredential) Credential() *model.Credential {
	return c.credential
}

func (m *mockData) AddCredential(c *model.Credential) (*model.Credential, error) {
	agent := m.agents[c.TenantID]

	n := model.NewCredential(c.TenantID, c)
	n.ID = faker.UUIDHyphenated()
	n.Created = time.Now().UTC()
	n.Cursor = model.TimeToCursor(&n.Created)
	for index := range n.Attributes {
		n.Attributes[index].ID = faker.UUIDHyphenated()
	}
	agent.credentials.append(newCredential(n))
	return n, nil
}

func (m *mockData) UpdateCredential(c *model.Credential) (*model.Credential, error) {
	agent := m.agents[c.TenantID]

	object := agent.credentials.objectForID(c.ID)
	if object == nil {
		return nil, errors.New("not found credential for id: " + c.ID)
	}
	updated := object.Copy()
	credential := updated.Credential()
	credential.Approved = c.Approved
	credential.Issued = c.Issued
	credential.Failed = c.Failed

	if !agent.credentials.replaceObjectForID(c.ID, updated) {
		return nil, errors.New("not found credential for id: " + c.ID)
	}
	return updated.Credential(), nil
}

func (m *mockData) GetCredential(id, tenantID string) (*model.Credential, error) {
	agent := m.agents[tenantID]

	c := agent.credentials.objectForID(id)
	if c == nil {
		return nil, errors.New("not found credential for id: " + id)
	}
	return c.Credential(), nil
}

func filterCredential(item apiObject) bool {
	c := item.Credential()
	return c.Issued != nil
}

func credentialConnectionFilter(id string) func(item apiObject) bool {
	return func(item apiObject) bool {
		c := item.Credential()
		if c.Issued != nil && c.ConnectionID == id {
			return true
		}
		return false
	}
}

func (m *mockItems) getCredentials(
	info *paginator.BatchInfo,
	filter func(item apiObject) bool,
) (connections *model.Credentials, err error) {
	state, hasNextPage, hasPreviousPage := m.credentials.getObjects(info, filter)
	res := make([]*model.Credential, len(state.objects))
	for i := range state.objects {
		res[i] = state.objects[i].Copy().Credential()
	}

	c := &model.Credentials{
		Credentials:     res,
		HasNextPage:     hasNextPage,
		HasPreviousPage: hasPreviousPage,
	}
	return c, nil
}

func (m *mockData) GetCredentials(
	info *paginator.BatchInfo,
	tenantID string,
	connectionID *string,
) (connections *model.Credentials, err error) {
	agent := m.agents[tenantID]

	if connectionID == nil {
		return agent.getCredentials(info, filterCredential)
	}
	return agent.getCredentials(info, credentialConnectionFilter(*connectionID))
}

func (m *mockData) GetCredentialCount(tenantID string, connectionID *string) (int, error) {
	agent := m.agents[tenantID]

	if connectionID == nil {
		return agent.credentials.count(filterCredential), nil
	}
	return agent.credentials.count(credentialConnectionFilter(*connectionID)), nil
}

func (m *mockData) GetConnectionForCredential(id, tenantID string) (*model.Connection, error) {
	credential, err := m.GetCredential(id, tenantID)
	if err != nil {
		return nil, err
	}
	return m.GetConnection(credential.ConnectionID, tenantID)
}
