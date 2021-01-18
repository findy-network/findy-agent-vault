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
	*base
	AgentID      string  `faker:"agentId"`
	Label        string  `faker:"first_name"`
	RawJWT       *string `faker:"-"`
	LastAccessed time.Time
}

func NewAgent(a *Agent) *Agent {
	if a != nil {
		return a.copy()
	}
	return &Agent{base: &base{}}
}

func (a *Agent) IsNewOnboard() bool {
	return a.Created == a.LastAccessed
}

func (a *Agent) copy() (n *Agent) {
	n = NewAgent(nil)
	if a.base != nil {
		n.base = a.base.copy()
	}
	n.AgentID = a.AgentID
	n.Label = a.Label
	n.RawJWT = a.RawJWT
	return n
}

func (a *Agent) ToNode() *model.User {
	return &model.User{
		ID:   a.ID,
		Name: a.Label,
	}
}
