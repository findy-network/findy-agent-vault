package model

import (
	"strconv"
	"time"
)

func timeToString(t *time.Time) string {
	return strconv.FormatInt(t.UnixNano()/time.Millisecond.Nanoseconds(), 10)
}

type base struct {
	ID      string `faker:"uuid_hyphenated"`
	Created time.Time
}

type Agent struct {
	*base
	AgentID      string `faker:"agentId"`
	Label        string `faker:"first_name"`
	LastAccessed time.Time
}

func NewAgent() *Agent { return &Agent{base: &base{}} }

func (a *Agent) Copy() (n *Agent) {
	n = NewAgent()
	n.AgentID = a.AgentID
	n.Label = a.Label
	return n
}
