package model

import (
	"math"
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

func (b *base) Copy() *base {
	baseCopy := *b
	return &baseCopy
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
	if a.base != nil {
		n.base = a.base.Copy()
	}
	n.AgentID = a.AgentID
	n.Label = a.Label
	return n
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
