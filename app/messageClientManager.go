package app

import (
	"fmt"
	"log"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
)

const (
	__INIT_MessageClientManager     int32 = 0
	__STARTING_MessageClientManager int32 = 1
	__STOPPED_MessageClientManager  int32 = 2
)

type MessageClientManager struct {
	clients map[MessageClient]string
	pipe    *MessagePipe

	validateClientID     ClientIDValidator
	onMessageClientClose func(MessageClient)
	logger               *log.Logger

	statusCode int32
}

func (manager *MessageClientManager) Join(client MessageClient) error {
	if __STOPPED_MessageClientManager == manager.statusCode {
		return ErrJoinClosedMessageClientManager
	}

	_, ok := manager.clients[client]
	if !ok {
		for i := 0; i < __MAX_GENERATING_CLIENT_ID_ATTEMPTS; i++ {
			id := manager.generateClientID()
			if manager.validateClientID(id) {
				manager.clients[client] = id

				client.setID(id)
				client.setStartAt(time.Now())
				client.setLogger(manager.logger)
				break
			}
		}

		_, ok = manager.clients[client]
		if !ok {
			return fmt.Errorf("generate client id failed")
		}

		// register close handler
		client.RegisterCloseHandler(manager.triggerMessageClientClose)

		if __STARTING_MessageClientManager == manager.statusCode {
			client.Start(manager.pipe)
		}
	}
	return nil
}

func (manager *MessageClientManager) Expel(client MessageClient) error {
	client.Stop()
	return client.Close()
}

func (manager *MessageClientManager) start() {
	if atomic.CompareAndSwapInt32(&manager.statusCode, __INIT_MessageClientManager, __STARTING_MessageClientManager) {
		for c := range manager.clients {
			c.Start(manager.pipe)
		}
	}
}

func (manager *MessageClientManager) stop() {
	if atomic.CompareAndSwapInt32(&manager.statusCode, __STARTING_MessageClientManager, __STOPPED_MessageClientManager) {
		// stopping all clients
		for c := range manager.clients {
			c.Stop()
		}

		// close all clients
		for c := range manager.clients {
			_ = c.Close()
		}
		return
	}

	atomic.StoreInt32(&manager.statusCode, __STOPPED_MessageClientManager)
}

func (manager *MessageClientManager) getClientID(client MessageClient) string {
	return manager.clients[client]
}

func (manager *MessageClientManager) generateClientID() string {
	id, err := uuid.NewUUID()
	if err != nil {
		panic(fmt.Sprintf("assert() cannot generate new uuid; %v", err))
	}
	return id.String()
}

func (manager *MessageClientManager) triggerMessageClientClose(client MessageClient) {
	if manager.onMessageClientClose != nil {
		manager.onMessageClientClose(client)
	}
	delete(manager.clients, client)
}
