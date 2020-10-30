package resolver

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"

	"github.com/golang/glog"

	"github.com/findy-network/findy-agent-vault/graph/model"
	"github.com/findy-network/findy-agent-vault/tools/data"
	"github.com/findy-network/findy-agent-vault/tools/faker"
	"github.com/lainio/err2"
)

var eventAddedObserver map[string]chan *model.EventEdge

func initEvents() {
	eventAddedObserver = map[string]chan *model.EventEdge{}
}

func (r *mutationResolver) MarkEventRead(ctx context.Context, input model.MarkReadInput) (node *model.Event, err error) {
	glog.V(logLevelMedium).Info("queryResolver:MarkEventRead, id: ", input.ID)

	state := data.State.Events
	node = state.MarkEventRead(input.ID)
	if node == nil {
		err = fmt.Errorf("event for id %s was not found", input.ID)
	}
	return
}

func (r *queryResolver) Events(
	_ context.Context,
	after *string, before *string,
	first, last *int) (c *model.EventConnection, err error) {
	defer err2.Return(&err)
	pagination := &PaginationParams{
		first:  first,
		last:   last,
		after:  after,
		before: before,
	}
	logPaginationRequest("queryResolver:events", pagination)

	state := data.State.Events
	afterIndex, beforeIndex, err := pick(state, pagination)
	err2.Check(err)

	return state.EventConnection(afterIndex, beforeIndex), nil
}

func (r *queryResolver) Event(ctx context.Context, id string) (node *model.Event, err error) {
	glog.V(logLevelMedium).Info("queryResolver:Event, id: ", id)

	state := data.State.Events
	node = state.EventForID(id)
	if node == nil {
		err = fmt.Errorf("event for id %s was not found", id)
	}
	return
}

func (r *subscriptionResolver) EventAdded(ctx context.Context) (<-chan *model.EventEdge, error) {
	id := "tenantId-" + strconv.FormatInt(time.Now().Unix(), 10)
	glog.V(logLevelMedium).Info("subscriptionResolver:EventAdded, id: ", id)

	// access user object: user := ctx.Value("user")

	events := make(chan *model.EventEdge, 1)

	go func() {
		<-ctx.Done()
		glog.V(logLevelMedium).Info("subscriptionResolver: event observer removed, id: ", id)
		delete(eventAddedObserver, id)
	}()

	eventAddedObserver[id] = events

	return events, nil
}

func doAddEvent(event *data.InternalEvent) {
	state := data.State.Events
	event.CreatedMs = time.Now().Unix()
	state.Append(event)
	glog.Infof("Added event %s", event.ID)
	for _, observer := range eventAddedObserver {
		observer <- event.ToEdge()
	}
}

func addEvent(description string, pType model.ProtocolType, pairwiseID string) {
	doAddEvent(&data.InternalEvent{
		ID:           uuid.New().String(),
		Read:         false,
		Description:  description,
		ProtocolType: pType,
		Type:         model.EventTypeNotification,
		PairwiseID:   pairwiseID,
	})
}

func (r *mutationResolver) AddRandomEvent(_ context.Context) (ok bool, err error) {
	glog.V(logLevelMedium).Info("mutationResolver:AddRandomEvent ")
	defer err2.Return(&err)

	events, err := faker.FakeEvents(1)
	err2.Check(err)

	doAddEvent(&events[0])
	ok = true

	return
}
