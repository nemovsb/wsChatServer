package app

const (
	EmptyEvent = iota
	DeleleConnectionEvent
)

type Event struct {
	Type int
	Data string
}

func NewEvent() *Event {
	return &Event{
		Type: EmptyEvent,
		Data: "",
	}
}
