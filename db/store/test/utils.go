package test

import (
	"github.com/findy-network/findy-agent-vault/db/fake"
	"github.com/findy-network/findy-agent-vault/db/model"
	"github.com/findy-network/findy-agent-vault/db/store"
)

func AddAgentAndConnections(db store.DB, agentID string, connectionCount int) (*model.Agent, []*model.Connection) {
	// add new agent
	input := model.NewAgent(nil)
	input.AgentID = agentID
	input.Label = agentID
	a, err := db.AddAgent(input)
	if err != nil {
		panic(err)
	}
	// add new connections
	connections := fake.AddConnections(db, a.ID, connectionCount)
	return a, connections
}
