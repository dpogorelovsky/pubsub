package pubsub

// Subscriber - minimum subscriber functionality
type Subscriber interface {
	Notify(interface{})
	Stop()
}
