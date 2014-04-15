package events

func NewTestEvent(value uint64) Event {
	return NewHeartbeat(value, 0, 0)
}

func GetTestEventType() Envelope_EventType {
	return Envelope_Heartbeat
}
