package data

import (
	"reflect"

	"github.com/findy-network/findy-agent-vault/tools/faker"

	"github.com/findy-network/findy-agent-vault/graph/model"

	our "github.com/findy-network/findy-agent-vault/tools/data/model"
)

type Data struct {
	connections *our.Items
	Messages    *our.Items
	Events      *our.Items
	Jobs        *our.Items
	User        *our.InternalUser
}

func InitState() *Data {
	state := &Data{
		connections: our.NewItems(reflect.TypeOf(model.Pairwise{}).Name()),
		Messages:    our.NewItems(reflect.TypeOf(model.BasicMessage{}).Name()),
		Events:      our.NewItems(reflect.TypeOf(model.Event{}).Name()),
		Jobs:        our.NewItems(reflect.TypeOf(model.Job{}).Name()),
	}
	state.User = faker.Run(state.connections, state.Events, state.Messages)
	state.sort()
	return state
}

func (state *Data) sort() {
	state.connections.Sort()
	state.Messages.Sort()
	state.Events.Sort()
}

func (state *Data) Connections() our.ConnectionItems {
	return state.connections.Connections()
}

func (state *Data) OutputForJob(id string) (output *model.JobOutput) {
	output = &model.JobOutput{
		Connection: nil,
		Message:    nil,
	}
	pType, pID := state.Jobs.JobProtocolForID(id)
	if pID != nil {
		switch pType {
		case model.ProtocolTypeConnection:
			output.Connection = state.Connections().PairwiseForID(*pID)
		case model.ProtocolTypeBasicMessage:
			output.Message = state.Messages.MessageForID(*pID)
		case model.ProtocolTypeNone:
		case model.ProtocolTypeCredential:
		case model.ProtocolTypeProof:
			break
		}
	}
	return
}
