package app

type MessageFormat uint

func (format MessageFormat) IsControl() bool {
	return format == CLOSE_MESSAGE ||
		format == PING_MESSAGE ||
		format == PONG_MESSAGE
}

func (format MessageFormat) IsData() bool {
	return format == TEXT_MESSAGE ||
		format == BINARY_MESSAGE
}
