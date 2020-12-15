package mock

import (
	"time"

	"github.com/bxcodec/faker/v3"
	"github.com/findy-network/findy-agent-vault/db/model"
)

func (m *mockData) AddAgent(a *model.Agent) (*model.Agent, error) {
	agent, ok := m.agentsByAgentID[a.AgentID]
	if !ok {
		agent = newState()
	}
	n := agent.agent

	now := time.Now().UTC()
	if !ok {
		n = model.NewAgent(a)
		n.ID = faker.UUIDHyphenated()
		n.Created = now
	} else {
		n = model.NewAgent(n)
	}
	n.LastAccessed = now
	agent.agent = n

	m.agents[n.ID] = agent
	m.agentsByAgentID[n.AgentID] = agent

	return n, nil
}

func (m *mockData) GetAgent(id, agentID *string) (*model.Agent, error) {
	var agent *mockItems
	if id != nil {
		agent = m.agents[*id]
	} else {
		agent = m.agentsByAgentID[*agentID]
	}
	return agent.agent, nil
}
