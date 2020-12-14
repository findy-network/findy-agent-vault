package mock

import (
	"reflect"

	"github.com/findy-network/findy-agent-vault/db/model"
	"github.com/findy-network/findy-agent-vault/db/store"
)

type mockItems struct {
	agent       *model.Agent
	connections *items
	credentials *items
	events      *items
}

func newState() *mockItems {
	state := &mockItems{
		agent:       nil,
		connections: newItems(reflect.TypeOf(model.Connection{}).Name()),
		credentials: newItems(reflect.TypeOf(model.Credential{}).Name()),
		events:      newItems(reflect.TypeOf(model.Event{}).Name()),
	}
	state.sort()
	return state
}

func InitState() store.DB {
	return newData()
}

func (m *mockItems) sort() {
	m.connections.sort()
	m.credentials.sort()
}

type apiObject interface {
	Identifier() string
	Created() uint64
	Connection() *model.Connection
	Credential() *model.Credential
	Event() *model.Event
	Copy() apiObject
}

type base struct{}

func (b *base) Connection() *model.Connection {
	panic("Object is not connection")
}

func (b *base) Credential() *model.Credential {
	panic("Object is not connection")
}

func (b *base) Event() *model.Event {
	panic("Object is not event")
}

type mockData struct {
	agents          map[string]*mockItems
	agentsByAgentID map[string]*mockItems
}

func newData() *mockData {
	return &mockData{
		agents:          make(map[string]*mockItems),
		agentsByAgentID: make(map[string]*mockItems),
	}
}

func (m *mockData) Close() {
	n := newData()
	m.agents = n.agents
	m.agentsByAgentID = n.agentsByAgentID
}
