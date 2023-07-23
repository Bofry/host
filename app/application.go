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

	sessionStateManager  SessionStateManager
	messageClientManager *MessageClientManager
	eventClient          EventClient

	messagePipe *MessagePipe
	eventPipe   *EventPipe

	// messageHandler MessageHandler
	// eventHandler   EventHandler
	errorHandler ErrorHandler

	messageChan chan *MessageSource
	eventChan   chan *Event
	errorChan   chan error
	worker      *ApplicationWorker

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

func (ap *Application) MessageClientManager() *MessageClientManager {
	return ap.messageClientManager
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

	ap.messageClientManager.start()
	ap.eventClient.Start(ap.eventPipe)
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

	ap.eventClient.Stop()
	ap.messageClientManager.stop()

	return ap.worker.stop(ctx)
}

func (ap *Application) alloc() {
	ap.messageChan = make(chan *MessageSource)
	ap.eventChan = make(chan *Event)
	ap.errorChan = make(chan error)

	ap.worker = &ApplicationWorker{
		logger:         ap.logger,
		receiveMessage: ap.acceptMessage,
		receiveEvent:   ap.receiveEvent,
		receiveError:   ap.receiveError,
		messageChan:    ap.messageChan,
		eventChan:      ap.eventChan,
		errorChan:      ap.errorChan,
	}
	ap.worker.alloc()

	ap.messagePipe = &MessagePipe{
		messageChan: ap.messageChan,
		errorChan:   ap.errorChan,
	}

	ap.eventPipe = &EventPipe{
		eventChan: ap.eventChan,
		errorChan: ap.errorChan,
	}

	ap.messageClientManager = &MessageClientManager{
		clients:              make(map[MessageClient]string),
		pipe:                 ap.messagePipe,
		validateClientID:     ap.validateClientID,
		onMessageClientClose: ap.triggerMessageClientClose,
	}
}

func (ap *Application) init() {
	if ap.initialized {
		return
	}

	defer func() {
		ap.initialized = true
	}()

	if ap.sessionStateManager == nil {
		ap.sessionStateManager = NewStdSessionStateManager()
	}

	ap.worker.init()
}

func (ap *Application) acceptMessage(source *MessageSource) {
	var (
		sessionID    = ap.messageClientManager.getClientID(source.Client)
		sessionState = ap.sessionStateManager.Load(sessionID)
	)

	if len(sessionID) == 0 {
		panic("assert() SessionID should be existed")
	}

	ctx := &Context{
		SessionID:             sessionID,
		SessionState:          sessionState,
		messageSender:         source.Client,
		eventForwarder:        ap.eventClient,
		logger:                ap.logger,
		invalidMessageHandler: nil, // be determined by MessageDispatcher
	}

	ap.worker.dispatchMessage(ctx, source.Message)
}

func (ap *Application) receiveEvent(event *Event) {
	ctx := &Context{
		SessionID:           "",
		SessionState:        nil,
		messageSender:       nil,
		eventForwarder:      ap.eventClient,
		logger:              ap.logger,
		invalidEventHandler: nil, // be determined by MessageDispatcher
	}

	ap.worker.dispatchEvent(ctx, event)
}

func (ap *Application) receiveError(err error) {
	if ap.errorHandler == nil {
		ap.logger.Println(err)
		return
	}
	ap.errorHandler(err)
}

func (ap *Application) configureProtocolResolver(resolver ProtocolResolver) {
	ap.worker.protocolResolver = resolver
}

func (ap *Application) configureInvalidMessageHandler(handler MessageHandler) {
	ap.worker.invalidMessageHandler = handler
}

func (ap *Application) configureInvalidEventHandler(handler EventHandler) {
	ap.worker.invalidEventHandler = handler
}

func (ap *Application) configureDefaultMessageHandler(handler MessageHandler) {
	ap.worker.defaultMessageHandler = handler
}

func (ap *Application) configureDefaultEventHandler(handler EventHandler) {
	ap.worker.defaultEventHandler = handler
}

func (ap *Application) configureMessageRouter(router MessageRouter) {
	ap.worker.messageRouter = router
}

func (ap *Application) configureEventRouter(router EventRouter) {
	ap.worker.eventRouter = router
}

func (ap *Application) validateClientID(id string) bool {
	return ap.sessionStateManager.TryCreate(id)
}

func (ap *Application) triggerMessageClientClose(client MessageClient) {
	var (
		id = ap.messageClientManager.getClientID(client)
	)
	ap.sessionStateManager.Delete(id)
}
