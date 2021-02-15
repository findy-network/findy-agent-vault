package mock

import (
	"time"

	"github.com/bxcodec/faker/v3"
	"github.com/findy-network/findy-agent-vault/db/model"
	"github.com/findy-network/findy-agent-vault/db/store"
	"github.com/findy-network/findy-agent-vault/paginator"
	"github.com/findy-network/findy-agent-vault/utils"
)

type mockMessage struct {
	*base
	message *model.Message
}

func (m *mockMessage) Created() uint64 {
	return model.TimeToCursor(&m.message.Created)
}

func (m *mockMessage) Identifier() string {
	return m.message.ID
}

func newMessage(m *model.Message) *mockMessage {
	var message *model.Message
	if m != nil {
		message = model.NewMessage(m.TenantID, m)
	}
	return &mockMessage{base: &base{}, message: message}
}

func (m *mockMessage) Copy() apiObject {
	return newMessage(m.message)
}

func (m *mockMessage) Message() *model.Message {
	return m.message
}

func (m *mockData) AddMessage(arg *model.Message) (*model.Message, error) {
	agent := m.agents.get(arg.TenantID)

	n := model.NewMessage(arg.TenantID, arg)
	n.ID = faker.UUIDHyphenated()
	n.Created = time.Now().UTC()
	n.Cursor = model.TimeToCursor(&n.Created)
	agent.messages.append(newMessage(n))
	return n, nil
}

func (m *mockData) UpdateMessage(arg *model.Message) (*model.Message, error) {
	agent := m.agents.get(arg.TenantID)

	object := agent.messages.objectForID(arg.ID)
	if object == nil {
		return nil, store.NewError(store.ErrCodeNotFound, "not found message for id: "+arg.ID)
	}
	updated := object.Copy()
	message := updated.Message()
	message.Delivered = arg.Delivered

	if !agent.messages.replaceObjectForID(arg.ID, updated) {
		panic("not found message for id: " + arg.ID)
	}
	return updated.Message(), nil
}

func (m *mockData) GetMessage(id, tenantID string) (*model.Message, error) {
	agent := m.agents.get(tenantID)

	msg := agent.messages.objectForID(id)
	if msg == nil {
		return nil, store.NewError(store.ErrCodeNotFound, "not found message for id: "+id)
	}
	return msg.Message(), nil
}

func messageConnectionFilter(id string) func(item apiObject) bool {
	return func(item apiObject) bool {
		return item.Message().ConnectionID == id
	}
}

func (m *mockItems) getMessages(
	info *paginator.BatchInfo,
	filter func(item apiObject) bool,
) (connections *model.Messages, err error) {
	state, hasNextPage, hasPreviousPage := m.messages.getObjects(info, filter)
	res := make([]*model.Message, len(state.objects))
	for i := range state.objects {
		res[i] = state.objects[i].Copy().Message()
	}

	c := &model.Messages{
		Messages:        res,
		HasNextPage:     hasNextPage,
		HasPreviousPage: hasPreviousPage,
	}
	return c, nil
}

func (m *mockData) GetMessages(
	info *paginator.BatchInfo,
	tenantID string,
	connectionID *string,
) (connections *model.Messages, err error) {
	agent := m.agents.get(tenantID)

	if connectionID == nil {
		return agent.getMessages(info, nil)
	}
	return agent.getMessages(info, messageConnectionFilter(*connectionID))
}

func (m *mockData) GetMessageCount(tenantID string, connectionID *string) (int, error) {
	agent := m.agents.get(tenantID)

	if connectionID == nil {
		return agent.messages.count(nil), nil
	}
	return agent.messages.count(messageConnectionFilter(*connectionID)), nil
}

func (m *mockData) GetConnectionForMessage(id, tenantID string) (*model.Connection, error) {
	message, err := m.GetMessage(id, tenantID)
	if err != nil {
		return nil, err
	}
	return m.GetConnection(message.ConnectionID, tenantID)
}

func (m *mockData) ArchiveMessage(id, tenantID string) error {
	agent := m.agents.get(tenantID)

	object := agent.messages.objectForID(id)
	if object == nil {
		return store.NewError(store.ErrCodeNotFound, "not found message for id: "+id)
	}

	now := utils.CurrentTime()

	n := model.NewMessage(tenantID, object.Message())
	n.Archived = &now

	if !agent.messages.replaceObjectForID(id, newMessage(n)) {
		panic("message not found")
	}

	return nil
}
