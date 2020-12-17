package mock

import (
	"sort"
	"sync"

	"github.com/findy-network/findy-agent-vault/paginator"
)

type items struct {
	objects []apiObject
	apiType string
	*sync.RWMutex
}

func newItems(apiType string) (i *items) {
	return &items{objects: make([]apiObject, 0), apiType: apiType, RWMutex: &sync.RWMutex{}}
}

func (i *items) append(object apiObject) {
	i.Lock()
	defer i.Unlock()
	i.objects = append(i.objects, object)
}

func (i *items) copy() *items {
	i.Lock()
	defer i.Unlock()
	n := newItems(i.apiType)
	for _, o := range i.objects {
		n.objects = append(n.objects, o.Copy())
	}
	return n
}

func (i *items) count(filter func(item apiObject) bool) (count int) {
	i.RLock()
	defer i.RUnlock()
	if filter == nil {
		count = len(i.objects)
	} else {
		count = len(i.filter(filterAndCopy(filter)).objects)
	}
	return
}

/*
func (i *items) randomID() *string {
	i.RLock()
	defer i.RUnlock()
	max := len(i.objects) - 1
	index := utils.Random(max)
	id := i.objects[index].Identifier()
	return &id
}*/

func (i *items) firstID() (id string) {
	i.RLock()
	defer i.RUnlock()
	if len(i.objects) > 0 {
		id = i.objects[0].Identifier()
	}
	return
}

func (i *items) lastID() (id string) {
	i.RLock()
	defer i.RUnlock()
	if len(i.objects) > 0 {
		id = i.objects[len(i.objects)-1].Identifier()
	}
	return
}

func (i *items) createdForIndex(index int) (created uint64) {
	i.RLock()
	defer i.RUnlock()
	created = i.objects[index].Created()
	return
}

func (i *items) sort() {
	i.Lock()
	defer i.Unlock()
	s := i.objects
	sort.Slice(s, func(i, j int) bool {
		return s[i].Created() < s[j].Created()
	})
}

func (i *items) filter(fn func(item apiObject) apiObject) *items {
	i.RLock()
	defer i.RUnlock()
	f := newItems(i.apiType)
	for index := range i.objects {
		res := fn(i.objects[index])
		if res != nil {
			f.append(res)
		}
	}
	return f
}

func (i *items) objectForID(id string) (o apiObject) {
	i.RLock()
	defer i.RUnlock()

	if id == "" {
		return
	}

	for _, item := range i.objects {
		if item.Identifier() == id {
			o = item.Copy()
			break
		}
	}

	return
}

func (i *items) replaceObjectForID(id string, o apiObject) (found bool) {
	i.RLock()
	defer i.RUnlock()

	if id == "" {
		return
	}

	for index := range i.objects {
		if i.objects[index].Identifier() == id {
			i.objects[index] = o
			found = true
			break
		}
	}

	return
}

func (i *items) getIndexes(info *paginator.BatchInfo) (afterIndex, beforeIndex int) {
	count := i.count(nil)
	before := info.Before
	after := info.After

	beforeIndex = count - 1
	if after != 0 || before != 0 {
		for index := 0; index < count; index++ {
			created := i.createdForIndex(index)
			if after > 0 && created <= after {
				nextIndex := index + 1
				if nextIndex < count {
					afterIndex = index + 1
				}
			}
			if before > 0 && created < before {
				beforeIndex = index
			}
			if (before > 0 && created > before) ||
				(before == 0 && created > after) {
				break
			}
		}
	}
	return afterIndex, beforeIndex
}

func filterAndCopy(filter func(item apiObject) bool) func(item apiObject) apiObject {
	return func(item apiObject) apiObject {
		if filter(item) {
			return item.Copy()
		}
		return nil
	}
}

func (i *items) getObjects(info *paginator.BatchInfo, filter func(item apiObject) bool) (state *items, hasNextPage, hasPreviousPage bool) {
	if filter != nil {
		state = i.filter(filterAndCopy(filter))
	} else {
		state = i.copy()
	}

	after, before := state.getIndexes(info)
	if !info.Tail {
		afterPlusFirst := after + (info.Count - 1)
		if before > afterPlusFirst {
			before = afterPlusFirst
		}
	} else {
		beforeMinusLast := before - (info.Count - 1)
		if after < beforeMinusLast {
			after = beforeMinusLast
		}
	}
	before++

	lastID := state.lastID()
	firstID := state.firstID()

	state.objects = state.objects[after:before]
	if len(state.objects) > 0 {
		hasNextPage = state.objects[len(state.objects)-1].Identifier() != lastID
		hasPreviousPage = state.objects[0].Identifier() != firstID
	}

	return state, hasNextPage, hasPreviousPage
}
