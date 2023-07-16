package app

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/Bofry/trace"
	"go.opentelemetry.io/otel/propagation"
)

type Application struct {
	Name string

	logger            *log.Logger
	tracerProvider    *trace.SeverityTracerProvider
	textMapPropagator propagation.TextMapPropagator

	messageSource MessageSource
	eventSource   EventSource

	messageHandler MessageHandler
	eventHandler   EventHandler
	errorHandler   ErrorHandler

	worker *Worker

	mutex       sync.Mutex
	initialized bool
	running     bool
	disposed    bool
}

func (ap *Application) TracerProvider() *trace.SeverityTracerProvider {
	return ap.tracerProvider
}

func (ap *Application) TextMapPropagator() propagation.TextMapPropagator {
	return ap.textMapPropagator
}

func (ap *Application) Start(ctx context.Context) error {
	if ap.disposed {
		return fmt.Errorf("the Application has been disposed")
	}
	if !ap.initialized {
		return fmt.Errorf("the Application havn't be initialized yet")
	}
	if ap.running {
		return nil
	}

	var err error
	ap.mutex.Lock()
	defer func() {
		if err != nil {
			ap.running = false
			ap.disposed = true
		}
		ap.mutex.Unlock()
	}()
	ap.running = true

	ap.messageSource.Receive(ap.worker.message)
	ap.eventSource.Notify(ap.worker.event)
	return ap.worker.start(ctx)
}

func (ap *Application) Stop(ctx context.Context) error {
	if ap.disposed {
		return nil
	}
	if !ap.running {
		return nil
	}

	ap.mutex.Lock()
	defer func() {
		ap.running = false
		ap.disposed = true
		ap.mutex.Unlock()
	}()

	ap.messageSource.Close()
	ap.eventSource.Close()

	return ap.worker.stop(ctx)
}

func (ap *Application) alloc() {
	ap.worker = &Worker{
		logger:         ap.logger,
		receiveMessage: ap.receiveMessage,
		receiveEvent:   ap.receiveEvent,
	}
}

func (ap *Application) init() {
	if ap.initialized {
		return
	}

	defer func() {
		ap.initialized = true
	}()

	ap.worker.init()
}

func (ap *Application) receiveMessage(message *Message) {
	ctx := &Context{
		logger:                ap.logger,
		invalidMessageHandler: nil, // be determined by MessageDispatcher
	}

	ap.worker.dispatchMessage(ctx, message)
}

func (ap *Application) receiveEvent(event *Event) {
	ctx := &Context{
		logger:              ap.logger,
		invalidEventHandler: nil, // be determined by MessageDispatcher
	}

	ap.worker.dispatchEvent(ctx, event)
}

func (ap *Application) receiveError(err error) {
	ap.errorHandler(err)
}

func (ap *Application) configureMessageCodeResolver(resolver MessageCodeResolver) {
	ap.worker.messageCodeResolver = resolver
}

func (ap *Application) configureInvalidMessageHandler(handler MessageHandler) {
	ap.worker.invalidMessageHandler = handler
}

func (ap *Application) configureInvalidEventHandler(handler EventHandler) {
	ap.worker.invalidEventHandler = handler
}

func (ap *Application) configureMessageRouter(router MessageRouter) {
	ap.worker.messageRouter = router
}

func (ap *Application) configureEventRouter(router EventRouter) {
	ap.worker.eventRouter = router
}
