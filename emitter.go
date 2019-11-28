package pubsub

import "log"

// Emitter - Event emitter interface
type Emitter interface {
	AddSubscriber(string, Subscriber)
	FireEvent(string, interface{})
	Stop()
}

// FiredEvent - fired event structure
type FiredEvent struct {
	name    string
	message interface{}
}

// SubscriberSignal - adding subscriber structure
type SubscriberSignal struct {
	eventName  string
	subscriber Subscriber
}

// NewEmitter - Emitter constructor
func NewEmitter(events []string) Emitter {
	e := &emitterInstance{
		stopSignal:          make(chan string),
		eventSignal:         make(chan FiredEvent),
		events:              make(map[string]bool),
		addSubscriberSignal: make(chan SubscriberSignal),
	}

	// setting possible/allowed events
	for _, val := range events {
		e.events[val] = true
	}

	go e.run()
	return e
}

type emitterInstance struct {
	events              map[string]bool
	subscribers         map[string][]Subscriber
	addSubscriberSignal chan SubscriberSignal
	eventSignal         chan FiredEvent
	stopSignal          chan string
}

// run - initializing emitter instance
// start listening to signals
func (e *emitterInstance) run() {
	for {
		select {
		case event := <-e.eventSignal:
			{
				e.handleFiredEvent(event)
			}
		case _ = <-e.stopSignal:
			{
				e.handleStop()
				return
			}
		}
	}
}

// AddSubscriber - adding subscriber to exact event
func (e *emitterInstance) AddSubscriber(eventName string, s Subscriber) {
	if _, ok := e.events[eventName]; !ok {
		panic("Adding event subscriber to non-existing event")
	}

	e.subscribers[eventName] = append(e.subscribers[eventName], s)
}

// FireEvent - we call it from outer code
func (e *emitterInstance) FireEvent(eventName string, message interface{}) {
	fired := FiredEvent{
		name:    eventName,
		message: message,
	}

	e.eventSignal <- fired
}

// Stop - sends to stopSignal channel
func (e *emitterInstance) Stop() {
	e.stopSignal <- "stop"
}

// handleStop - just a wrapper for some logic, not to have it in select{}
// when stopping emitter we close listening 'stop' and 'event' channels
// also call Stop() handler on every subscriber
func (e *emitterInstance) handleStop() {
	log.Println("stopping emitter...")
	for _, event := range e.subscribers {
		for _, sub := range event {
			sub.Stop()
		}
	}
	close(e.eventSignal)
	close(e.stopSignal)
}

// handleFiredEvent - wrapping logic that calls notify() on all subscribers
func (e *emitterInstance) handleFiredEvent(event FiredEvent) {
	if _, ok := e.events[event.name]; !ok {
		panic("Firing non-existing event")
	}

	for _, sub := range e.subscribers[event.name] {
		sub.Notify(event.message)
	}
}
