package resolver

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/golang/glog"

	"github.com/findy-network/findy-agent-api/graph/model"
	"github.com/findy-network/findy-agent-api/tools/data"
	"github.com/findy-network/findy-agent-api/tools/faker"
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

	user := ctx.Value("user")
	fmt.Printf("%v", user)

	events := make(chan *model.EventEdge, 1)

	go func() {
		<-ctx.Done()
		glog.V(logLevelMedium).Info("subscriptionResolver: event observer removed, id: ", id)
		delete(eventAddedObserver, id)
	}()

	eventAddedObserver[id] = events

	return events, nil
}

func (r *mutationResolver) AddRandomEvent(_ context.Context) (bool, error) {
	glog.V(logLevelMedium).Info("mutationResolver:AddRandomEvent ")

	state := data.State.Events
	events, err := faker.FakeEvents(1)
	if err == nil {
		event := &events[0]
		event.CreatedMs = time.Now().Unix()
		state.Append(event)
		glog.Infof("Added random event %s", events[0].ID)
		for _, observer := range eventAddedObserver {
			observer <- events[0].ToEdge()
		}
	}
	return true, err
}
