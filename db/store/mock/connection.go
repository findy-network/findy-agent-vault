package mock

import (
	"errors"
	"time"

	"github.com/bxcodec/faker/v3"
	"github.com/findy-network/findy-agent-vault/db/model"
	"github.com/findy-network/findy-agent-vault/paginator"
)

type mockConnection struct {
	connection *model.Connection
}

func (c *mockConnection) Created() uint64 {
	return model.TimeToCursor(&c.connection.Created)
}

func (c *mockConnection) Identifier() string {
	return c.connection.ID
}

func (c *mockConnection) Copy() apiObject {
	return &mockConnection{c.connection.Copy()}
}

func (c *mockConnection) Connection() *model.Connection {
	return c.connection
}

func (c *mockConnection) Credential() *model.Credential {
	panic("Object is not credential")
}

func (m *mockData) AddConnection(c *model.Connection) (*model.Connection, error) {
	agent := m.agents[c.TenantID]

	n := c.Copy()
	n.ID = faker.UUIDHyphenated()
	n.Created = time.Now().UTC()
	n.Cursor = model.TimeToCursor(&n.Created)
	object := &mockConnection{n}
	agent.connections.append(object)

	// generate different timestamps for items
	time.Sleep(time.Millisecond)

	return n, nil
}

func (m *mockData) GetConnection(id string, tenantID string) (*model.Connection, error) {
	agent := m.agents[tenantID]

	c := agent.connections.objectForID(id)
	if c == nil {
		return nil, errors.New("not found connection for id: " + id)
	}
	return c.Connection(), nil
}

func (m *mockData) GetConnections(info *paginator.BatchInfo, tenantID string) (connections *model.Connections, err error) {
	agent := m.agents[tenantID]

	state, hasNextPage, hasPreviousPage := agent.connections.getObjects(info, nil)
	res := make([]*model.Connection, len(state.objects))
	for i := range state.objects {
		res[i] = state.objects[i].Copy().Connection()
	}

	c := &model.Connections{
		Connections:     res,
		HasNextPage:     hasNextPage,
		HasPreviousPage: hasPreviousPage,
	}
	return c, nil
}

func (m *mockData) GetConnectionCount(tenantID string) (int, error) {
	agent := m.agents[tenantID]

	return agent.connections.count(nil), nil
}
