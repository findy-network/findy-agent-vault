package update

import (
	"context"
	"strconv"
	"sync"

	dbModel "github.com/findy-network/findy-agent-vault/db/model"
	"github.com/findy-network/findy-agent-vault/graph/model"
	"github.com/findy-network/findy-agent-vault/utils"
	"github.com/golang/glog"
	"github.com/lainio/err2"
	"github.com/lainio/err2/try"
)

type subscription struct {
	channel  chan *model.EventEdge
	tenantID string
}

type subscriberRegister struct {
	*sync.RWMutex
	subscriptions map[string]*subscription
	agents        map[string][]string
}

func newSubscriberRegister() *subscriberRegister {
	return &subscriberRegister{
		RWMutex:       &sync.RWMutex{},
		subscriptions: make(map[string]*subscription),
		agents:        make(map[string][]string),
	}
}

func (s *subscriberRegister) notify(tenantID string, event *dbModel.Event) {
	s.RLock()
	defer s.RUnlock()

	agentSubscriptions, ok := s.agents[tenantID]
	if !ok {
		utils.LogTrace().Infof("Skipping notifications, no subscriptions for %s", tenantID)
		return
	}

	for _, subscriptionID := range agentSubscriptions {
		subscription, ok := s.subscriptions[subscriptionID]
		if !ok {
			glog.Errorf("No subscription channel found for subscription ID %s", subscriptionID)
			continue
		}

		eventEdge := event.ToEdge()
		utils.LogMed().Infof("Sending event %v to tenant %s", eventEdge, tenantID)
		subscription.channel <- eventEdge
	}
}

func (s *subscriberRegister) add(tenantID string) (subscriptionID string, eventChannel <-chan *model.EventEdge) {
	s.Lock()
	defer s.Unlock()

	utils.LogMed().Infof("Add subscription for tenant %s", tenantID)

	const timeLen = 10
	subscriptionID = tenantID + "-" + strconv.FormatInt(utils.CurrentTimeMs(), timeLen)
	newSubscription := &subscription{
		tenantID: tenantID,
		channel:  make(chan *model.EventEdge, 1),
	}
	eventChannel = newSubscription.channel
	s.subscriptions[subscriptionID] = newSubscription

	subscriptions, ok := s.agents[tenantID]
	if !ok {
		subscriptions = make([]string, 0)
	}
	subscriptions = append(subscriptions, subscriptionID)
	s.agents[tenantID] = subscriptions

	return
}

func (s *subscriberRegister) remove(subscriptionID string) {
	s.Lock()
	defer s.Unlock()

	subscription, ok := s.subscriptions[subscriptionID]
	if !ok {
		glog.Errorf("Attempted to remove non-existing subscription with ID %s", subscriptionID)
		return
	}
	tenantID := subscription.tenantID
	delete(s.subscriptions, subscriptionID)

	subscriptions, ok := s.agents[tenantID]
	if !ok {
		glog.Errorf("Attempted to remove non-existing agent subscription with ID %s", subscriptionID)
		return
	}
	index := -1
	for i := 0; i < len(subscriptions); i++ {
		if subscriptions[i] == subscriptionID {
			index = i
			break
		}
	}
	if index < 0 {
		glog.Errorf("Subscription with ID %s was not found in agent %s register", subscriptionID, tenantID)
		return
	}
	if len(subscriptions) > 1 {
		subscriptions[index] = subscriptions[0]
		subscriptions = subscriptions[1:]
	} else {
		subscriptions = make([]string, 0)
	}
	s.agents[tenantID] = subscriptions

	utils.LogMed().Infof("Subscription %s was removed for tenant %s", subscriptionID, tenantID)
}

func (r *Updater) EventAdded(ctx context.Context) (ch <-chan *model.EventEdge, err error) {
	defer err2.Return(&err)

	tenant := try.To1(r.GetAgent(ctx))

	id, events := r.eventSubscribers.add(tenant.ID)
	utils.LogMed().Info("subscriptionResolver:EventAdded, id: ", id)

	go func() {
		<-ctx.Done()
		utils.LogMed().Info("subscriptionResolver: event observer removed, id: ", id)
		r.eventSubscribers.remove(id)
	}()

	return events, err
}
