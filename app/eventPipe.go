package app

type EventPipe struct {
	eventChan chan *Event
	errorChan chan error
}

func (pipe *EventPipe) Forward(event *Event) {
	pipe.eventChan <- event
}

func (pipe *EventPipe) Error(err error) {
	pipe.errorChan <- err
}
