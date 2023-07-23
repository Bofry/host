package app

import "reflect"

const (
	__MAX_GENERATING_CLIENT_ID_ATTEMPTS int = 64

	__MODULE_APP_FIELD                   = "App"
	__APP_DEFAULT_MESSAGE_HANDLER_METHOD = "DefaultMessageHandler"
	__APP_DEFAULT_EVENT_HANDLER_METHOD   = "DefaultEventHandler"
	__APP_EVENT_CLIENT_FIELD             = "EventClient"

	__LOGGER_PREFIX_FORMAT = "[host/app/%s] "

	InvalidChannel string = "?"

	TAG_PROTOCOL       = "protocol"
	TAG_CHANNEL        = "channel"
	TAG_OPT_EXPAND_ENV = "@ExpandEnv"
	OPT_ON             = "on"
	OPT_OFF            = "off"
)

var (
	typeOfMessageHandler = reflect.TypeOf(MessageHandler(nil))
	typeOfEventHandler   = reflect.TypeOf(EventHandler(nil))
	typeOfEventClient    = reflect.TypeOf((*EventClient)(nil)).Elem()

	StandardProtocolMessageFormats = []MessageFormat{
		TEXT_MESSAGE, BINARY_MESSAGE,
	}
)

const (
	UNKNOWN_MESSAGE MessageFormat = 0
	TEXT_MESSAGE    MessageFormat = 1
	BINARY_MESSAGE  MessageFormat = 2
	CLOSE_MESSAGE   MessageFormat = 8
)

const (
	APP            TargetValueRole = "APP"
	MODULE_OPTIONS TargetValueRole = "MODULE_OPTIONS"
)

type (
	MessageFormat   uint
	TargetValueRole string
)

type (
	SessionStateManager interface {
		Load(id string) SessionState
		Update(id string, state SessionState)
		Delete(id string)
		TryCreate(id string) bool
	}

	SessionState interface {
		CanVisit() bool
		Visit(func(k, v interface{}))
		Value(k interface{}) interface{}
		SetValue(k, v interface{})
		Lock()
		Unlock()
	}

	MessageSender interface {
		Send(format MessageFormat, payload []byte) error
	}

	MessageClient interface {
		Start(*MessagePipe)
		Stop()
		Close() error
		RegisterCloseHandler(func(MessageClient))

		MessageSender
	}

	EventForwarder interface {
		Forward(channel string, payload []byte) error
	}

	EventClient interface {
		Start(*EventPipe)
		Stop()
		Close() error

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

	ClientIDValidator func(string) bool

	ProtocolResolver func(format MessageFormat, payload []byte) string

	Module interface {
		ModuleOptions() []ApplicationBuildingOption
	}

	ApplicationBuildingOption interface {
		apply(*Application) error
	}

	ModuleBindingOption interface {
		apply(reflect.Value, TargetValueRole) error
	}

	StandardProtocol interface {
		ConfigureProtocol(registry *StandardProtocolRegistry)
		ReplyMessage(sender MessageSender)
	}

	MessageHasher func(Message) uint64
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
