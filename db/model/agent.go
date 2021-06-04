package model

import (
	"time"

	"github.com/findy-network/findy-agent-vault/graph/model"
)

type Agents struct {
	Agents          []*Agent
	HasNextPage     bool
	HasPreviousPage bool
}

type Agent struct {
	base
	AgentID      string `faker:"agentId"`
	Label        string `faker:"first_name"`
	RawJWT       string `faker:"-"`
	LastAccessed time.Time
}

func (a *Agent) IsNewOnboard() bool {
	return a.Created == a.LastAccessed
}

func (a *Agent) ToNode() *model.User {
	return &model.User{
		ID:   a.ID,
		Name: a.Label,
	}
}
