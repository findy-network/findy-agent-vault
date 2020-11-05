package data

import (
	"reflect"
	"sort"

	"github.com/findy-network/findy-agent-vault/graph/model"

	our "github.com/findy-network/findy-agent-vault/tools/data/model"
)

type Data struct {
	Connections *our.Items
	Events      *our.Items
	Jobs        *our.Items
	User        *our.InternalUser
}

func InitState(scratch bool) *Data {
	state := &Data{
		Connections: our.NewItems(reflect.TypeOf(model.Pairwise{}).Name()),
		Events:      our.NewItems(reflect.TypeOf(model.Event{}).Name()),
		Jobs:        our.NewItems(reflect.TypeOf(model.Job{}).Name()),
		User:        &user,
	}
	state.initStateAndSort(scratch)
	return state
}

func (state *Data) initStateAndSort(scratch bool) {
	sort.Slice(connections, func(i, j int) bool {
		return connections[i].Created() < connections[j].Created()
	})

	sort.Slice(events, func(i, j int) bool {
		return events[i].Created() < events[j].Created()
	})

	if !scratch {
		for index := range connections {
			state.Connections.Append(&connections[index])
		}
		state.Connections.Sort()

		for index := range events {
			state.Events.Append(&events[index])
		}
		state.Events.Sort()
	}
}

func (state *Data) MarkEventRead(id string) *model.Event {
	if state.Events.MarkEventRead(id) {
		return state.Events.EventForID(id, state.Connections)
	}

	return nil
}
