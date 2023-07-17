package app

import (
	"context"
	"log"
	"sync"
	"time"
)

var (
	_ context.Context = new(Context)
)

type Context struct {
	SessionID    string
	SessionState SessionState

	messageSender  MessageSender
	eventForwarder EventForwarder

	context context.Context
	logger  *log.Logger

	invalidMessageHandler MessageHandler
	invalidEventHandler   EventHandler

	values     map[interface{}]interface{}
	valuesOnce sync.Once
}

// Deadline implements context.Context.
func (*Context) Deadline() (deadline time.Time, ok bool) {
	return
}

// Done implements context.Context.
func (*Context) Done() <-chan struct{} {
	return nil
}

// Err implements context.Context.
func (*Context) Err() error {
	return nil
}

// Value implements context.Context.
func (ctx *Context) Value(key any) any {
	if key == nil {
		return nil
	}
	if ctx.values != nil {
		v := ctx.values[key]
		if v != nil {
			return v
		}
	}
	if ctx.context != nil {
		return ctx.context.Value(key)
	}
	return nil
}

func (ctx *Context) SetValue(key interface{}, value interface{}) {
	if key == nil {
		return
	}
	if ctx.values == nil {
		ctx.valuesOnce.Do(func() {
			if ctx.values == nil {
				ctx.values = make(map[interface{}]interface{})
			}
		})
	}
	ctx.values[key] = value
}

func (ctx *Context) Forward(channel string, payload []byte) {
	if ctx.eventForwarder == nil {
		return
	}
	ctx.eventForwarder.Forward(channel, payload)
}

func (ctx *Context) Send(format MessageFormat, payload []byte) {
	if ctx.messageSender == nil {
		return
	}
	ctx.messageSender.Send(format, payload)
}

func (ctx *Context) Logger() *log.Logger {
	return ctx.logger
}

func (ctx *Context) InvalidMessage(message *Message) {
	if ctx.invalidMessageHandler != nil {
		ctx.invalidMessageHandler(ctx, message)
	}
}

func (ctx *Context) InvalidEvent(event *Event) {
	if ctx.invalidMessageHandler != nil {
		ctx.invalidEventHandler(ctx, event)
	}
}
