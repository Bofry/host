package app

const (
	__LOGGER_PREFIX_FORMAT = "[host/app/%s] "

	InvalidChannel string = "?"
)

const (
	UNKNOWN MessageFormat = 0
	TEXT    MessageFormat = 1
	BINARY  MessageFormat = 2
	CLOSE   MessageFormat = 8
)

type (
	MessageFormat int
)

type (
	MessageSender interface {
		Send(format MessageFormat, payload []byte) error
	}

	MessageSource interface {
		Start(chan *Message, chan error)
		Stop() error

		MessageSender
	}

	EventForwarder interface {
		Forward(channel string, payload []byte) error
	}

	EventSource interface {
		Start(chan *Event, chan error)
		Stop() error

		EventForwarder
	}

	EventDelegate interface {
		OnAck(event *Event)
		OnRetry(event *Event)
		OnAbort(event *Event)
	}

	MessageContent interface {
		Decode(format MessageFormat, body []byte) error
		Encode() (MessageFormat, []byte)
		Validate() error
	}

	MessageHandler func(ctx *Context, message *Message)
	EventHandler   func(ctx *Context, event *Event) error
	ErrorHandler   func(err error)

	MessageCodeResolver func(format MessageFormat, payload []byte) string

	ModuleOptionCollection []ApplicationBuildingOption

	ApplicationBuildingOption interface {
		apply(*Application) error
	}
)

/*
MODE 1 - REQ/REP and notify
C -> S :wg.Add(1) -> MessageHandler            // wait  1
S -> Q :ctx.Forward()
C <- S :ctx.Send() -> wg.Done()                // wait  0

MODE 3 and more complex conditions
C -> S :wg.Add(1) -> MessageHandler            // wait: 1
S -> Q :ctx.Forward()
C <- S :ctx.Send() -> wg.Done()                // wait  0
----------------------------------------------------------------
Q -> S :wg.Add(1) -> EventHandler              // wait  1
S -> Q :ctx.Forward()
C <- S :ctx.Send() -> wg.Done()                // wait  0

MODE 2 - SUB/REP
Q -> S :wg.Add(1) -> EventHandler              // wait: 1
S -> Q :ctx.Forward()
C <- S :ctx.Send() -> wg.Done()                // wait  0
*/
