package app

import (
	"context"
	"log"
	"sync"
)

type Worker struct {
	logger *log.Logger

	receiveMessage func(*MessageSource)
	receiveEvent   func(*Event)
	receiveError   func(error)

	messageCodeResolver   MessageCodeResolver
	invalidMessageHandler MessageHandler
	invalidEventHandler   EventHandler

	messageRouter MessageRouter
	eventRouter   EventRouter

	messageChan chan *MessageSource
	eventChan   chan *Event
	errorChan   chan error
	done        chan struct{}

	wg    sync.WaitGroup
	mutex sync.Mutex

	initialized bool
}

func (w *Worker) alloc() {
	w.done = make(chan struct{})
}

func (w *Worker) init() {
	if w.initialized {
		return
	}

	w.mutex.Lock()
	if w.initialized {
		return
	}

	defer func() {
		w.initialized = true
		w.mutex.Unlock()
	}()

	if w.receiveMessage == nil {
		w.receiveMessage = func(m *MessageSource) {}
	}
	if w.receiveEvent == nil {
		w.receiveEvent = func(e *Event) {}
	}
}

func (w *Worker) start(ctx context.Context) error {
	var (
		message = w.messageChan
		event   = w.eventChan
		error   = w.errorChan
		done    = w.done
	)

	var kontinue bool = true
	go func() {
		for kontinue || len(message) > 0 || len(event) > 0 {
			select {
			case v, ok := <-message:
				if ok {
					w.wg.Add(1)
					defer func() {
						w.wg.Done()
					}()

					w.receiveMessage(v)
				}
			case v, ok := <-event:
				if ok {
					w.wg.Add(1)
					defer w.wg.Done()

					w.receiveEvent(v)
				}
			case v, ok := <-error:
				if ok {
					w.wg.Add(1)
					defer w.wg.Done()

					w.receiveError(v)
				}
			case <-done:
				w.logger.Println("Stopping")
				kontinue = false
				break
			}
		}
	}()
	return nil
}

func (w *Worker) stop(ctx context.Context) error {
	close(w.done)
	w.wg.Wait()
	return nil
}

func (w *Worker) dispatchMessage(ctx *Context, message *Message) {
	var (
		router = w.messageRouter
	)

	ctx.invalidMessageHandler = w.invalidMessageHandler

	if w.messageCodeResolver != nil {
		code := w.messageCodeResolver(message.Format, message.Body)

		handler := router.Get(code)
		if handler != nil {
			handler(ctx, message)
			return
		}
	}
	ctx.InvalidMessage(message)
}

func (w *Worker) dispatchEvent(ctx *Context, event *Event) {
	var (
		router = w.eventRouter
	)

	handler := router.Get(event.Channel)
	if handler != nil {
		handler(ctx, event)
		return
	}
	ctx.InvalidEvent(event)
}
