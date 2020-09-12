package app

import (
	"sync"
)

type EventSystem struct {
	mu        sync.RWMutex
	listeners []chan []byte
}

func NewEventSystem() *EventSystem {
	return &EventSystem{
		listeners: make([]chan []byte, 0),
	}
}

func (es *EventSystem) Send(msg []byte) {
	es.mu.RLock()
	defer es.mu.RUnlock()

	for _, l := range es.listeners {
		l <- msg
	}
}

func (es *EventSystem) Listen(listener chan []byte) {
	es.mu.Lock()
	defer es.mu.Unlock()

	es.listeners = append(es.listeners, listener)
}

func (es *EventSystem) Close() {
	es.mu.Lock()
	defer es.mu.Unlock()

	for _, l := range es.listeners {
		close(l)
	}
}
