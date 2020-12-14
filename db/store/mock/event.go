package mock

import (
	"errors"
	"time"

	"github.com/bxcodec/faker/v3"
	"github.com/findy-network/findy-agent-vault/db/model"
	"github.com/findy-network/findy-agent-vault/paginator"
)

type mockEvent struct {
	*base
	event *model.Event
}

func (e *mockEvent) Created() uint64 {
	return model.TimeToCursor(&e.event.Created)
}

func (e *mockEvent) Identifier() string {
	return e.event.ID
}

func newEvent(e *model.Event) *mockEvent {
	var event *model.Event
	if e != nil {
		event = model.NewEvent(e)
	}
	return &mockEvent{base: &base{}, event: event}
}

func (e *mockEvent) Copy() apiObject {
	return newEvent(e.event)
}

func (e *mockEvent) Event() *model.Event {
	return e.event
}

func (m *mockData) AddEvent(e *model.Event) (*model.Event, error) {
	agent := m.agents[e.TenantID]

	n := model.NewEvent(e)
	n.ID = faker.UUIDHyphenated()
	n.Created = time.Now().UTC()
	n.Cursor = model.TimeToCursor(&n.Created)
	agent.events.append(newEvent(n))

	// generate different timestamps for items
	time.Sleep(time.Millisecond)

	return n, nil
}

func (m *mockData) MarkEventRead(tenantID, id string) (*model.Event, error) {
	agent := m.agents[tenantID]

	object := agent.events.objectForID(id)
	if object == nil {
		return nil, errors.New("not found event for id: " + id)
	}
	updated := object.Copy()
	event := updated.Event()
	event.Read = true

	if !agent.events.replaceObjectForID(id, updated) {
		return nil, errors.New("not found event for id: " + id)
	}
	return updated.Event(), nil
}

func (m *mockData) GetEvent(id, tenantID string) (*model.Event, error) {
	agent := m.agents[tenantID]

	e := agent.events.objectForID(id)
	if e == nil {
		return nil, errors.New("not found event for id: " + id)
	}
	return e.Event(), nil
}

func (m *mockItems) getEvents(
	info *paginator.BatchInfo,
	filter func(item apiObject) bool,
) (connections *model.Events, err error) {
	state, hasNextPage, hasPreviousPage := m.events.getObjects(info, filter)
	res := make([]*model.Event, len(state.objects))
	for i := range state.objects {
		res[i] = state.objects[i].Copy().Event()
	}

	c := &model.Events{
		Events:          res,
		HasNextPage:     hasNextPage,
		HasPreviousPage: hasPreviousPage,
	}
	return c, nil
}

func (m *mockData) GetEvents(info *paginator.BatchInfo, tenantID string) (connections *model.Events, err error) {
	agent := m.agents[tenantID]

	return agent.getEvents(info, nil)
}

func (m *mockData) GetEventCount(tenantID string) (int, error) {
	agent := m.agents[tenantID]

	return agent.events.count(nil), nil
}

func eventConnectionFilter(id string) func(item apiObject) bool {
	return func(item apiObject) bool {
		e := item.Event()
		if e.ConnectionID != nil && *e.ConnectionID == id {
			return true
		}
		return false
	}
}

func (m *mockData) GetConnectionEvents(
	info *paginator.BatchInfo,
	tenantID,
	connectionID string,
) (connections *model.Events, err error) {
	agent := m.agents[tenantID]
	return agent.getEvents(info, eventConnectionFilter(connectionID))
}

func (m *mockData) GetConnectionEventCount(tenantID, connectionID string) (int, error) {
	agent := m.agents[tenantID]
	return agent.events.count(eventConnectionFilter(connectionID)), nil
}
