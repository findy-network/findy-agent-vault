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
}

func newState() *mockItems {
	state := &mockItems{
		connections: newItems(reflect.TypeOf(model.Connection{}).Name()),
		credentials: newItems(reflect.TypeOf(model.Credential{}).Name()),
		agent:       nil,
	}
	state.sort()
	return state
}

func InitState() store.DB {
	return newData()
}

func (state *mockItems) sort() {
	state.connections.sort()
	state.credentials.sort()
}

type apiObject interface {
	Identifier() string
	Created() uint64
	Connection() *model.Connection
	Credential() *model.Credential
	Copy() apiObject
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

func (d *mockData) Close() {
	d = newData()
}
