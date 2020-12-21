package mock

import (
	"reflect"
	"time"

	"github.com/bxcodec/faker/v3"
	"github.com/findy-network/findy-agent-vault/db/model"
	"github.com/findy-network/findy-agent-vault/paginator"
)

type mockAgent struct {
	*base
	agent *model.Agent
}

func (a *mockAgent) Created() uint64 {
	return model.TimeToCursor(&a.agent.Created)
}

func (a *mockAgent) Identifier() string {
	return a.agent.ID
}

func newAgent(a *model.Agent) *mockAgent {
	var agent *model.Agent
	if a != nil {
		agent = model.NewAgent(a)
	}
	return &mockAgent{base: &base{}, agent: agent}
}

func (a *mockAgent) Copy() apiObject {
	return newAgent(a.agent)
}

func (a *mockAgent) Agent() *model.Agent {
	return a.agent
}

func (m *mockData) GetListenerAgents(info *paginator.BatchInfo) (*model.Agents, error) {
	agentsData := m.agents.getAgents()

	agents := newItems(reflect.TypeOf(model.Agent{}).Name())
	for _, a := range agentsData {
		agents.append(newAgent(a))
	}
	agents.sort()

	state, hasNextPage, hasPreviousPage := agents.getObjects(info, func(item apiObject) bool {
		return item.Agent().RawJWT != ""
	})
	res := make([]*model.Agent, len(state.objects))
	for i := range state.objects {
		res[i] = state.objects[i].Copy().Agent()
	}

	a := &model.Agents{
		Agents:          res,
		HasNextPage:     hasNextPage,
		HasPreviousPage: hasPreviousPage,
	}
	return a, nil
}

func (m *mockData) AddAgent(a *model.Agent) (*model.Agent, error) {
	agent := m.agents.getByAgentID(a.AgentID)
	add := false
	if agent == nil {
		agent = newState()
		add = true
	}
	n := model.NewAgent(agent.agent)

	now := time.Now().UTC()
	if add {
		n = model.NewAgent(a)
		n.ID = faker.UUIDHyphenated()
		n.Created = now
	} else {
		n.AgentID = a.AgentID
		n.Label = a.Label
		n.RawJWT = a.RawJWT
	}
	n.LastAccessed = now
	n.TenantID = n.ID
	agent.agent = n

	m.agents.set(n.ID, n.AgentID, agent)

	return n, nil
}

func (m *mockData) GetAgent(id, agentID *string) (*model.Agent, error) {
	var agent *mockItems
	if id != nil {
		agent = m.agents.get(*id)
	} else {
		agent = m.agents.getByAgentID(*agentID)
	}
	return agent.agent, nil
}
