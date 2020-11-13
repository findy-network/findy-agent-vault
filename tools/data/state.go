package data

import (
	"reflect"

	"github.com/findy-network/findy-agent-vault/tools/faker"

	"github.com/findy-network/findy-agent-vault/graph/model"

	our "github.com/findy-network/findy-agent-vault/tools/data/model"
)

type Data struct {
	Connections *our.Items
	Messages    *our.Items
	Events      *our.Items
	Jobs        *our.Items
	User        *our.InternalUser
}

func InitState() *Data {
	state := &Data{
		Connections: our.NewItems(reflect.TypeOf(model.Pairwise{}).Name()),
		Messages:    our.NewItems(reflect.TypeOf(model.BasicMessage{}).Name()),
		Events:      our.NewItems(reflect.TypeOf(model.Event{}).Name()),
		Jobs:        our.NewItems(reflect.TypeOf(model.Job{}).Name()),
	}
	state.User = faker.Run(state.Connections, state.Events, state.Messages)
	state.sort()
	return state
}

func (state *Data) sort() {
	state.Connections.Sort()
	state.Messages.Sort()
	state.Events.Sort()
}

func (state *Data) MarkEventRead(id string) *model.EventEdge {
	if state.Events.MarkEventRead(id) {
		return state.Events.EventForID(id)
	}

	return nil
}
