package app

type MessagePipe struct {
	messageChan chan *MessageSource
	errorChan   chan error
}

func (pipe *MessagePipe) Forward(client MessageClient, message *Message) {
	pipe.messageChan <- &MessageSource{
		Message: message,
		Client:  client,
	}
}

func (pipe *MessagePipe) Error(err error) {
	pipe.errorChan <- err
}
