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
	proofs      *items
	events      *items
	messages    *items
}

func newState() *mockItems {
	state := &mockItems{
		agent:       nil,
		connections: newItems(reflect.TypeOf(model.Connection{}).Name()),
		credentials: newItems(reflect.TypeOf(model.Credential{}).Name()),
		proofs:      newItems(reflect.TypeOf(model.Proof{}).Name()),
		events:      newItems(reflect.TypeOf(model.Event{}).Name()),
		messages:    newItems(reflect.TypeOf(model.Message{}).Name()),
	}
	return state
}

func InitState() store.DB {
	return newData()
}

func (m *mockItems) sort() {
	m.connections.sort()
	m.credentials.sort()
	m.events.sort()
	m.messages.sort()
}

type apiObject interface {
	Identifier() string
	Created() uint64
	Connection() *model.Connection
	Credential() *model.Credential
	Proof() *model.Proof
	Event() *model.Event
	Message() *model.Message
	Copy() apiObject
}

type base struct{}

func (b *base) Connection() *model.Connection {
	panic("Object is not connection")
}

func (b *base) Credential() *model.Credential {
	panic("Object is not connection")
}

func (b *base) Proof() *model.Proof {
	panic("Object is not proof")
}

func (b *base) Event() *model.Event {
	panic("Object is not event")
}

func (b *base) Message() *model.Message {
	panic("Object is not message")
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
