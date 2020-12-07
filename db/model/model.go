package model

import (
	"encoding/base64"
	"reflect"
	"strconv"
	"time"
)

func CreateCursor(created uint64, object interface{}) string {
	typeName := reflect.TypeOf(object).Name()
	return base64.StdEncoding.EncodeToString(
		[]byte(typeName + ":" + strconv.FormatUint(created, 10)),
	)
}

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
