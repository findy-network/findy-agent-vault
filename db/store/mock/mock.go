package mock

import (
	"reflect"
	"sync"

	"github.com/findy-network/findy-agent-vault/db/model"
	"github.com/findy-network/findy-agent-vault/db/store"
)

type mockItems struct {
	agent       *model.Agent
	connections *items
	credentials *items
	proofs      *items
	messages    *items
	events      *items
	jobs        *items
}

func newState() *mockItems {
	state := &mockItems{
		agent:       nil,
		connections: newItems(reflect.TypeOf(model.Connection{}).Name()),
		credentials: newItems(reflect.TypeOf(model.Credential{}).Name()),
		proofs:      newItems(reflect.TypeOf(model.Proof{}).Name()),
		messages:    newItems(reflect.TypeOf(model.Message{}).Name()),
		events:      newItems(reflect.TypeOf(model.Event{}).Name()),
		jobs:        newItems(reflect.TypeOf(model.Job{}).Name()),
	}
	return state
}

func InitState() store.DB {
	return newData()
}

func (m *mockItems) sort() {
	m.connections.sort()
	m.credentials.sort()
	m.proofs.sort()
	m.messages.sort()
	m.events.sort()
	m.jobs.sort()
}

type apiObject interface {
	Identifier() string
	Created() uint64
	Copy() apiObject
	Connection() *model.Connection
	Credential() *model.Credential
	Proof() *model.Proof
	Message() *model.Message
	Event() *model.Event
	Job() *model.Job
	Agent() *model.Agent
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

func (b *base) Message() *model.Message {
	panic("Object is not message")
}

func (b *base) Event() *model.Event {
	panic("Object is not event")
}

func (b *base) Job() *model.Job {
	panic("Object is not job")
}

func (b *base) Agent() *model.Agent {
	panic("Object is not agent")
}

type mockData struct {
	agents *agentLedger
}

func newData() *mockData {
	return &mockData{
		agents: &agentLedger{
			RWMutex:         &sync.RWMutex{},
			agents:          make(map[string]*mockItems),
			agentsByAgentID: make(map[string]*mockItems),
		},
	}
}

func (m *mockData) Close() {
	n := newData()
	m.agents = n.agents
}

type agentLedger struct {
	*sync.RWMutex
	agents          map[string]*mockItems
	agentsByAgentID map[string]*mockItems
}

func (a *agentLedger) get(tenantID string) *mockItems {
	a.RLock()
	defer a.RUnlock()

	items := a.agents[tenantID]

	return items
}

func (a *agentLedger) getByAgentID(agentID string) *mockItems {
	a.RLock()
	defer a.RUnlock()

	items := a.agentsByAgentID[agentID]

	return items
}

func (a *agentLedger) set(tenantID, agentID string, items *mockItems) {
	a.Lock()
	defer a.Unlock()
	a.agents[tenantID] = items
	a.agentsByAgentID[agentID] = items
}

func (a *agentLedger) getAgents() []*model.Agent {
	a.Lock()
	defer a.Unlock()

	agents := make([]*model.Agent, 0)
	for _, agentData := range a.agents {
		agents = append(agents, agentData.agent)
	}
	return agents
}
