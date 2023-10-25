package app

var (
	_ MessageClient = new(RestrictedMessageClient)
)

type RestrictedMessageClient struct {
	client MessageClient

	RestrictedMessageClientInfo
}

func NewRestrictedMessageClient(client MessageClient) *RestrictedMessageClient {
	if v, ok := client.(*RestrictedMessageClient); ok {
		return v
	}

	return &RestrictedMessageClient{
		client:                      client,
		RestrictedMessageClientInfo: RestrictedMessageClientInfo(client),
	}
}

// Close implements MessageClient.
func (*RestrictedMessageClient) Close() error {
	panic("the operation Close() is restricted")
}

// RegisterCloseHandler implements MessageClient.
func (*RestrictedMessageClient) RegisterCloseHandler(func(MessageClient)) {
	panic("the operation RegisterCloseHandler() is restricted")
}

// Send implements MessageClient.
func (*RestrictedMessageClient) Send(format MessageFormat, payload []byte) error {
	panic("the operation Send() is restricted")
}

// Start implements MessageClient.
func (*RestrictedMessageClient) Start(*MessagePipe) {
	panic("the operation Start() is restricted")
}

// Stop implements MessageClient.
func (*RestrictedMessageClient) Stop() {
	panic("the operation Stop() is restricted")
}
