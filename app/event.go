package app

import "sync/atomic"

type Event struct {
	Channel  string
	Payload  []byte
	Delegate EventDelegate

	responded int32
}

func (e *Event) Ack() {
	if !atomic.CompareAndSwapInt32(&e.responded, 0, 1) {
		return
	}
	e.Delegate.OnAck(e)
}

func (e *Event) Retry() {
	if !atomic.CompareAndSwapInt32(&e.responded, 0, 1) {
		return
	}
	e.Delegate.OnRetry(e)
}

func (e *Event) Abort() {
	if !atomic.CompareAndSwapInt32(&e.responded, 0, 1) {
		return
	}
	e.Delegate.OnAbort(e)
}
