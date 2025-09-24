package events

type Event interface {
	Handlers() []any
	UID() uint16
}
