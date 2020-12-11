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

	n := a.Copy()
	if n.ID == "" {
		n.ID = faker.UUIDHyphenated()
		n.Created = time.Now().UTC()
		n.LastAccessed = time.Now().UTC()
	} else {
		n.LastAccessed = time.Now().UTC()
	}
	agent.agent = n

	if !ok {
		m.agents[n.ID] = agent
		m.agentsByAgentID[n.AgentID] = agent
	}

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
