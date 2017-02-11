package events

import (
	"sync"
)

// Event emitter
type Event struct {
	sync.Mutex
	listeners []*Listener
}

// Listener instance
type Listener struct {
	fn     func(...interface{})
	once   bool
	parent *Event
}

// New Event emitter
func New() *Event {
	return &Event{listeners: []*Listener{}}
}

// Drop Event emitter
func (e *Event) Drop() {
	e = nil
}

// Remove listener
func (l *Listener) Remove() {
	l.parent.RemoveListener(l)
}

func (e *Event) addListener(fn func(...interface{}), once bool) (listener *Listener) {
	listener = &Listener{fn: fn, once: once, parent: e}
	e.Lock()
	defer e.Unlock()
	e.listeners = append(e.listeners, listener)
	return
}

// On - create a new listener
func (e *Event) On(fn func(...interface{})) (listener *Listener) {
	listener = e.addListener(fn, false)
	return
}

// Once - create a new one-time listener
func (e *Event) Once(fn func(...interface{})) (listener *Listener) {
	listener = e.addListener(fn, true)
	return
}

// RemoveListener - remove event's listener
func (e *Event) RemoveListener(l *Listener) *Event {
	listeners := []*Listener{}
	e.Lock()
	defer e.Unlock()
	for _, v := range e.listeners {
		if v != l {
			listeners = append(listeners, v)
		}
	}
	e.listeners = listeners
	return e
}

// Clear removes all listeners from (all/event)
func (e *Event) Clear() *Event {
	e.Lock()
	defer e.Unlock()
	e.listeners = []*Listener{}
	return e
}

// ListenersCount returns the count of listeners in the speicifed event
func (e *Event) ListenersCount() int {
	e.Lock()
	defer e.Unlock()
	return len(e.listeners)
}

// Emit new event
func (e *Event) Emit(args ...interface{}) *Event {
	listeners := []*Listener{}
	e.Lock()
	defer e.Unlock()
	for _, v := range e.listeners {
		v.fn(args...)
		if !v.once {
			listeners = append(listeners, v)
		}
	}
	e.listeners = listeners
	return e
}
