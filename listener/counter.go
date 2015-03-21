package listener

import (
	"sync"
)

type tagCounter struct {
	data map[string]int
	sync.RWMutex
}

func newCounter() *tagCounter {
	return &tagCounter{
		data: make(map[string]int),
	}
}

func (t *tagCounter) incCount(tag string, value int) {
	if value == 0 {
		return
	}

	t.Lock()
	defer t.Unlock()

	if _, ok := t.data[tag]; ok {
		t.data[tag] += value
	} else {
		t.data[tag] = value
	}
}

func (t *tagCounter) getDataAndClear() (counts map[string]int) {
	t.Lock()
	defer t.Unlock()

	counts = make(map[string]int)

	for key, value := range t.data {
		counts[key] = value
		delete(t.data, key)
	}

	return
}
