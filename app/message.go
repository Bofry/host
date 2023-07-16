package app

type Message struct {
	Format MessageFormat
	Body   []byte
}

func (m *Message) DecodeContent(v MessageContent) error {
	err := v.Decode(m.Format, m.Body)
	if err != nil {
		return err
	}
	return v.Validate()
}
