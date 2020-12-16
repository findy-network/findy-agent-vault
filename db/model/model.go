package model

import (
	"math"
	"strconv"
	"time"

	"github.com/findy-network/findy-agent-vault/graph/model"
)

func timeToString(t *time.Time) string {
	return strconv.FormatInt(t.UnixNano()/time.Millisecond.Nanoseconds(), 10)
}

type base struct {
	ID       string `faker:"uuid_hyphenated"`
	TenantID string
	Cursor   uint64
	Created  time.Time
}

func (b *base) copy() *base {
	baseCopy := *b
	return &baseCopy
}

type Agent struct {
	*base
	AgentID      string `faker:"agentId"`
	Label        string `faker:"first_name"`
	LastAccessed time.Time
}

func NewAgent(a *Agent) *Agent {
	if a != nil {
		return a.copy()
	}
	return &Agent{base: &base{}}
}

func (a *Agent) copy() (n *Agent) {
	n = NewAgent(nil)
	if a.base != nil {
		n.base = a.base.copy()
	}
	n.AgentID = a.AgentID
	n.Label = a.Label
	return n
}

func (a *Agent) ToNode() *model.User {
	return &model.User{
		ID:   a.ID,
		Name: a.Label,
	}
}

func copyTime(t *time.Time) *time.Time {
	var res *time.Time
	if t != nil {
		ts := *t
		res = &ts
	}
	return res
}

func TimeToCursor(t *time.Time) uint64 {
	return uint64(math.Round(float64(t.UnixNano()) / float64(time.Millisecond.Nanoseconds())))
}
