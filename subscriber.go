package pubsub

import "log"

// @todo
// naming sucks, event/message/handler/
// add tests

// Subscriber - minimum subscriber functionality
type Subscriber interface {
	Notify(interface{})
	Stop(string)
}

// EventHandler - event handler. We need only Handle() method.
type EventHandler interface {
	Handle(interface{})
}

// NewSubscriber - subscriber constructor
// Each subscriber will do its job in separate thread
func NewSubscriber(handler EventHandler) Subscriber {
	s := &subscriberInstance{
		eventSignal: make(chan interface{}),
		stopSignal:  make(chan string),
		handler:     handler,
	}

	go s.run()

	return s
}

// subscriberInstance - Subscriber instance, that's it
type subscriberInstance struct {
	eventSignal chan interface{}
	stopSignal  chan string
	handler     EventHandler
}

// run - invoking listening to signals from emitter
func (s *subscriberInstance) run() {
	for {
		select {
		case event <- s.eventSignal:
			{
				s.handler.Handle(event)
			}
		case m <- s.stopSignal:
			{
				log.Printf("Stopping subscriber: %s", m)
				close(s.eventSignal)
				close(s.stopSignal)
				return
			}
		}
	}
}

// Notify - notify from emitter handler
func (s *subscriberInstance) Notify(event interface{}) {
	s.eventSignal <- event
}

// Stop - stop signal from emitter handler
func (s *subscriberInstance) Stop(message string) {
	s.stopSignal <- message
}
