package app

import (
	"log"
	"time"
)

var (
	_ MessageClientInfoImpl = new(MessageClientInfo)
)

type MessageClientInfo struct {
	id      string
	startAt time.Time
	logger  *log.Logger
}

func NewMessageClientInfo() *MessageClientInfo {
	return new(MessageClientInfo)
}

// ID implements MessageClientInfoImpl.
func (info *MessageClientInfo) ID() string {
	return info.id
}

// Logger implements MessageClientInfoImpl.
func (info *MessageClientInfo) Logger() *log.Logger {
	return info.logger
}

// StartAt implements MessageClientInfoImpl.
func (info *MessageClientInfo) StartAt() time.Time {
	return info.startAt
}

// setID implements MessageClientInfoImpl.
func (info *MessageClientInfo) setID(v string) {
	info.id = v
}

// setStartAt implements MessageClientInfoImpl.
func (info *MessageClientInfo) setStartAt(v time.Time) {
	info.startAt = v
}

// Logger implements MessageClientInfoImpl.
func (info *MessageClientInfo) setLogger(v *log.Logger) {
	info.logger = v
}
