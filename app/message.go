package app

import "bytes"

type Message struct {
	Format   MessageFormat
	Protocol string
	Body     []byte
}

func (m *Message) DecodeContent(v MessageContent) error {
	err := v.Decode(m.Format, m.Body)
	if err != nil {
		return err
	}
	return v.Validate()
}

func (m Message) Equals(other Message) bool {
	return (m.Format == other.Format) &&
		(m.Protocol == other.Protocol) &&
		bytes.Equal(m.Body, other.Body)
}
