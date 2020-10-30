package data

import (
	"reflect"
	"sort"
	"sync"

	"github.com/findy-network/findy-agent-vault/tools/utils"

	"github.com/findy-network/findy-agent-vault/graph/model"
)

type Items struct {
	items   []APIObject
	apiType string
	mutex   sync.RWMutex
}

func (i *Items) Append(object APIObject) {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	i.items = append(i.items, object)
}

func (i *Items) Count() (count int) {
	i.mutex.RLock()
	defer i.mutex.RUnlock()
	count = len(i.items)
	return
}

func (i *Items) RandomID() (id string) {
	i.mutex.RLock()
	defer i.mutex.RUnlock()
	max := len(i.items) - 1
	index := utils.Random(max)
	id = i.items[index].Identifier()
	return
}

func (i *Items) FirstID() (id string) {
	i.mutex.RLock()
	defer i.mutex.RUnlock()
	id = i.items[0].Identifier()
	return
}

func (i *Items) LastID() (id string) {
	i.mutex.RLock()
	defer i.mutex.RUnlock()
	id = i.items[len(i.items)-1].Identifier()
	return
}

func (i *Items) CreatedForIndex(index int) (created int64) {
	i.mutex.RLock()
	defer i.mutex.RUnlock()
	created = i.items[index].Created()
	return
}

func (i *Items) MinCreated() (created int64) {
	i.mutex.RLock()
	defer i.mutex.RUnlock()
	created = i.items[0].Created()
	return
}

func (i *Items) MaxCreated() (created int64) {
	i.mutex.RLock()
	defer i.mutex.RUnlock()
	created = i.items[len(i.items)-1].Created()
	return
}

func (i *Items) Sort() {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	s := i.items
	sort.Slice(s, func(i, j int) bool {
		return s[i].Created() < s[j].Created()
	})
}

func (i *Items) PairwiseForID(id string) *model.Pairwise {
	var node *model.Pairwise

	i.mutex.RLock()
	defer i.mutex.RUnlock()

	for _, item := range i.items {
		if item.Identifier() == id {
			node = item.Pairwise().ToNode()
			break
		}
	}

	return node
}

func (i *Items) PairwiseConnection(after, before int) *model.PairwiseConnection {
	i.mutex.RLock()
	result := i.items[after:before]
	totalCount := len(result)

	edges := make([]*model.PairwiseEdge, totalCount)
	nodes := make([]*model.Pairwise, totalCount)
	for index, pairwise := range result {
		node := pairwise.Pairwise().ToNode()
		edges[index] = &model.PairwiseEdge{
			Cursor: CreateCursor(pairwise.Pairwise().CreatedMs, model.Pairwise{}),
			Node:   node,
		}
		nodes[index] = node
	}
	i.mutex.RUnlock()

	var startCursor, endCursor *string
	if totalCount > 0 {
		startCursor = &edges[0].Cursor
		endCursor = &edges[totalCount-1].Cursor
	}
	p := &model.PairwiseConnection{
		Edges: edges,
		Nodes: nodes,
		PageInfo: &model.PageInfo{
			EndCursor:       endCursor,
			HasNextPage:     edges[len(edges)-1].Node.ID != i.LastID(),
			HasPreviousPage: edges[0].Node.ID != i.FirstID(),
			StartCursor:     startCursor,
		},
		TotalCount: totalCount,
	}
	return p
}

func (i *Items) EventForID(id string) *model.Event {
	var node *model.Event

	i.mutex.RLock()
	defer i.mutex.RUnlock()

	for _, item := range i.items {
		if item.Identifier() == id {
			node = item.Event().ToNode()
			break
		}
	}

	return node
}

func (i *Items) EventConnection(after, before int) *model.EventConnection {
	i.mutex.RLock()
	result := i.items[after:before]
	totalCount := len(result)

	edges := make([]*model.EventEdge, totalCount)
	nodes := make([]*model.Event, totalCount)
	for index, event := range result {
		node := event.Event().ToNode()
		edges[index] = &model.EventEdge{
			Cursor: CreateCursor(event.Event().CreatedMs, model.Event{}),
			Node:   node,
		}
		nodes[index] = node
	}
	i.mutex.RUnlock()

	var startCursor, endCursor *string
	if totalCount > 0 {
		startCursor = &edges[0].Cursor
		endCursor = &edges[totalCount-1].Cursor
	}
	c := &model.EventConnection{
		Edges: edges,
		Nodes: nodes,
		PageInfo: &model.PageInfo{
			EndCursor:       endCursor,
			HasNextPage:     edges[len(edges)-1].Node.ID != i.LastID(),
			HasPreviousPage: edges[0].Node.ID != i.FirstID(),
			StartCursor:     startCursor,
		},
		TotalCount: totalCount,
	}
	return c
}

func (i *Items) MarkEventRead(id string) *model.Event {
	i.mutex.Lock()
	defer i.mutex.Unlock()

	for _, item := range i.items {
		if item.Identifier() == id {
			event := item.Event()
			event.Read = true
			return event.ToNode()
		}
	}
	return nil
}

type Data struct {
	Connections *Items
	Events      *Items
	User        *InternalUser
}

var State = &Data{
	Connections: &Items{items: make([]APIObject, 0), apiType: reflect.TypeOf(model.Pairwise{}).Name()},
	Events:      &Items{items: make([]APIObject, 0), apiType: reflect.TypeOf(model.Event{}).Name()},
	User:        &user,
}

func InitState() {
	InitStateAndSort(false)
}

func InitStateAndSort(scratch bool) {
	sort.Slice(connections, func(i, j int) bool {
		return connections[i].Created() < connections[j].Created()
	})

	sort.Slice(events, func(i, j int) bool {
		return events[i].Created() < events[j].Created()
	})

	if !scratch {
		for index := range connections {
			State.Connections.items = append(State.Connections.items, &connections[index])
		}
		State.Connections.Sort()

		for index := range events {
			State.Events.items = append(State.Events.items, &events[index])
		}
		State.Events.Sort()
	}
}
