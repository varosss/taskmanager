package store

import (
	"sort"
)

const INITIAL_ID = 1

type MemoryStore[T any] struct {
	data   map[int]T
	nextId int
}

func NewMemoryStore[T any]() *MemoryStore[T] {
	return &MemoryStore[T]{
		data:   make(map[int]T),
		nextId: 1,
	}
}

func (ms *MemoryStore[T]) SetData(data map[int]T) {
	ms.data = data
}

func (ms *MemoryStore[T]) SetNextId(nextId int) {
	ms.nextId = nextId
}

func (ms *MemoryStore[T]) Add(items []T, setId func(*T, int)) {
	for _, item := range items {
		setId(&item, ms.nextId)
		ms.data[ms.nextId] = item
		ms.nextId++
	}
}

func (ms *MemoryStore[T]) Delete(id int) {
	delete(ms.data, id)
}

func (ms *MemoryStore[T]) Update(items []T, getId func(T) int) {
	for _, item := range items {
		ms.data[getId(item)] = item
	}
}

func (ms *MemoryStore[T]) Get(id int) (T, bool) {
	item, ok := ms.data[id]
	return item, ok
}

func (ms *MemoryStore[T]) Data() map[int]T {
	return ms.data
}

func (ms *MemoryStore[T]) List() []T {
	keys := make([]int, 0, len(ms.data))
	for k := range ms.data {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	items := make([]T, 0, len(ms.data))
	for _, k := range keys {
		items = append(items, ms.data[k])
	}

	return items
}
