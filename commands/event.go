package commands

type Event struct {
	Event string
	Data  string
	Id    string
	Retry int
}

type EventSender interface {
	SendEvent(Event)
}
