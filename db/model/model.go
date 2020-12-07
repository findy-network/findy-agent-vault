package model

import (
	"time"
)

type base struct {
	ID      string `faker:"uuid_hyphenated"`
	Created time.Time
}

type Agent struct {
	*base
	AgentID string
	Label   string `faker:"first_name"`
}

func NewAgent() *Agent { return &Agent{base: &base{}} }

func (a *Agent) Copy() (n *Agent) {
	n = NewAgent()
	n.AgentID = a.AgentID
	n.Label = a.Label
	return n
}
