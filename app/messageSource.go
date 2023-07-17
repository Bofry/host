package app

type MessageSource struct {
	Message *Message
	Client  MessageClient
}

type MessageWorker struct {
}
