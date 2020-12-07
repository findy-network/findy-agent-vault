package model

import (
	"sort"
	"sync"

	utils "github.com/findy-network/findy-agent-vault/tools/tools"
)

type Items struct {
	items   []APIObject
	apiType string
	mutex   sync.RWMutex
}

func NewItems(apiType string) (i *Items) {
	return &Items{items: make([]APIObject, 0), apiType: apiType}
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

func (i *Items) RandomID() *string {
	i.mutex.RLock()
	defer i.mutex.RUnlock()
	max := len(i.items) - 1
	index := utils.Random(max)
	id := i.items[index].Identifier()
	return &id
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

func (i *Items) Filter(fn func(item APIObject) APIObject) *Items {
	i.mutex.RLock()
	defer i.mutex.RUnlock()
	f := NewItems(i.apiType)
	for index := range i.items {
		res := fn(i.items[index])
		if res != nil {
			f.Append(res)
		}
	}
	return f
}
