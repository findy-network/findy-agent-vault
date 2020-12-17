package mock

import (
	"time"

	"github.com/bxcodec/faker/v3"
	"github.com/findy-network/findy-agent-vault/db/model"
)

func (m *mockData) AddAgent(a *model.Agent) (*model.Agent, error) {
	agent := m.agents.getByAgentID(a.AgentID)
	add := false
	if agent == nil {
		agent = newState()
		add = true
	}
	n := agent.agent

	now := time.Now().UTC()
	if add {
		n = model.NewAgent(a)
		n.ID = faker.UUIDHyphenated()
		n.Created = now
	} else {
		n = model.NewAgent(n)
	}
	n.LastAccessed = now
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
