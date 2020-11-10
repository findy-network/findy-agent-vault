package data

import (
	our "github.com/findy-network/findy-agent-vault/tools/data/model"
)

type Pairwise = our.InternalPairwise
type Event = our.InternalEvent
type User = our.InternalUser

var connections = []Pairwise{}

var events = []Event{}

var user = User{}
