package mock

import (
	"errors"
	"time"

	"github.com/bxcodec/faker/v3"
	"github.com/findy-network/findy-agent-vault/db/model"
	"github.com/findy-network/findy-agent-vault/paginator"
)

type mockCredential struct {
	credential *model.Credential
}

func (c *mockCredential) Created() uint64 {
	return model.TimeToCursor(&c.credential.Created)
}

func (c *mockCredential) Identifier() string {
	return c.credential.ID
}

func (c *mockCredential) Copy() apiObject {
	return &mockCredential{c.credential.Copy()}
}

func (c *mockCredential) Connection() *model.Connection {
	panic("Object is not connection")
}

func (c *mockCredential) Credential() *model.Credential {
	return c.credential
}

func (m *mockData) AddCredential(c *model.Credential) (*model.Credential, error) {
	agent := m.agents[c.TenantID]

	n := c.Copy()
	n.ID = faker.UUIDHyphenated()
	n.Created = time.Now().UTC()
	n.Cursor = model.TimeToCursor(&n.Created)
	for index := range n.Attributes {
		n.Attributes[index].ID = faker.UUIDHyphenated()
	}
	agent.credentials.append(&mockCredential{n})

	// generate different timestamps for items
	time.Sleep(time.Millisecond)

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

func (m *mockData) GetCredentials(info *paginator.BatchInfo, tenantID string) (connections *model.Credentials, err error) {
	agent := m.agents[tenantID]

	return agent.getCredentials(info, filterCredential)
}

func (m *mockData) GetCredentialCount(tenantID string) (int, error) {
	agent := m.agents[tenantID]

	return agent.credentials.count(filterCredential), nil
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

func (m *mockData) GetConnectionCredentials(
	info *paginator.BatchInfo,
	tenantID,
	connectionID string,
) (connections *model.Credentials, err error) {
	agent := m.agents[tenantID]
	return agent.getCredentials(info, credentialConnectionFilter(connectionID))
}

func (m *mockData) GetConnectionCredentialCount(tenantID, connectionID string) (int, error) {
	agent := m.agents[tenantID]
	return agent.credentials.count(credentialConnectionFilter(connectionID)), nil
}
