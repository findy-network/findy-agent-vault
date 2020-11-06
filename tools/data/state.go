package data

import (
	"reflect"
	"sort"

	"github.com/findy-network/findy-agent-vault/tools/faker"

	"github.com/findy-network/findy-agent-vault/graph/model"

	our "github.com/findy-network/findy-agent-vault/tools/data/model"
)

type Data struct {
	Connections *our.Items
	Events      *our.Items
	Jobs        *our.Items
	User        *our.InternalUser
}

func InitState() *Data {
	state := &Data{
		Connections: our.NewItems(reflect.TypeOf(model.Pairwise{}).Name()),
		Events:      our.NewItems(reflect.TypeOf(model.Event{}).Name()),
		Jobs:        our.NewItems(reflect.TypeOf(model.Job{}).Name()),
		User:        &user,
	}
	faker.Run(state.Connections, state.Events)
	state.sort()
	return state
}

func (state *Data) sort() {
	sort.Slice(connections, func(i, j int) bool {
		return connections[i].Created() < connections[j].Created()
	})

	sort.Slice(events, func(i, j int) bool {
		return events[i].Created() < events[j].Created()
	})
	state.Connections.Sort()
	state.Events.Sort()
}

func (state *Data) MarkEventRead(id string) *model.Event {
	if state.Events.MarkEventRead(id) {
		return state.Events.EventForID(id, state.Connections, state.Jobs)
	}

	return nil
}
