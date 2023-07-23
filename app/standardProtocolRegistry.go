package app

import (
	"encoding/binary"
	"hash"
	"hash/fnv"
	"sync"
)

type StandardProtocolRegistry struct {
	keyValues map[uint64][]*StandardProtocolEntity

	messageHashers []MessageHasher
	mutax          sync.Mutex
}

func NewProtocolMessageRegistry() *StandardProtocolRegistry {
	instance := &StandardProtocolRegistry{
		keyValues: make(map[uint64][]*StandardProtocolEntity),
	}

	instance.messageHashers = []MessageHasher{
		instance.computeFNV1,
		instance.computeFNV1a,
	}
	return instance
}

func (r *StandardProtocolRegistry) Add(message Message, protocol StandardProtocol) {
	if protocol != nil {
		// exists?
		for _, computeHashCode := range r.messageHashers {
			hashcode := computeHashCode(message)
			kvs, ok := r.keyValues[hashcode]
			if ok {
				for _, elem := range kvs {
					if elem.Message.Equals(message) {
						elem.Protocol = protocol
						return
					}
				}
			}
		}

		// does not exists
		hashcode := r.computeFNV1(message)
		// add eneity
		r.mutax.Lock()
		{
			values := r.keyValues[hashcode]
			values = append(values, &StandardProtocolEntity{
				Message:  message,
				Protocol: protocol,
			})
			r.keyValues[hashcode] = values
		}
		r.mutax.Unlock()
	}
}

func (r *StandardProtocolRegistry) Get(message Message) StandardProtocol {
	for _, computeHashCode := range r.messageHashers {
		hashcode := computeHashCode(message)
		kvs, ok := r.keyValues[hashcode]
		if ok {
			for _, elem := range kvs {
				if elem.Message.Equals(message) {
					return elem.Protocol
				}
			}
		}
	}
	return nil
}

func (r *StandardProtocolRegistry) Remove(message Message) {
	for _, computeHashCode := range r.messageHashers {
		hashcode := computeHashCode(message)
		values, ok := r.keyValues[hashcode]
		if ok {
			for i, elem := range values {
				if elem.Message.Equals(message) {
					r.mutax.Lock()
					// remove eneity
					values = append(values[:i], values[i+1:]...)
					if len(values) == 0 {
						delete(r.keyValues, hashcode)
					} else {
						r.keyValues[hashcode] = values
					}
					r.mutax.Unlock()
					return
				}
			}
		}
	}
}

func (r *StandardProtocolRegistry) Visit(visitor func(Message, StandardProtocol)) {
	for _, values := range r.keyValues {
		for _, pair := range values {
			visitor(pair.Message, pair.Protocol)
		}
	}
}

func (r *StandardProtocolRegistry) computeHash(hasher hash.Hash64, message Message) uint64 {
	var (
		formatBytes = make([]byte, 8)
	)
	binary.LittleEndian.PutUint32(formatBytes, uint32(message.Format))

	hasher.Write(formatBytes)
	hasher.Write(message.Body)
	return hasher.Sum64()
}

func (r *StandardProtocolRegistry) computeFNV1(message Message) uint64 {
	return r.computeHash(fnv.New64(), message)
}

func (r *StandardProtocolRegistry) computeFNV1a(message Message) uint64 {
	return r.computeHash(fnv.New64a(), message)
}
