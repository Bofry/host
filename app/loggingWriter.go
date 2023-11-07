package app

import (
	"io"
	"runtime/debug"
)

var _ io.Writer = new(LoggingWriter)

type LoggingWriter struct {
	writer io.Writer

	EnableStackTrace bool
}

func NewLoggingWriter(writer io.Writer) *LoggingWriter {
	if writer == nil {
		writer = io.Discard
	}

	return &LoggingWriter{
		writer: writer,
	}
}

// Write implements io.Writer.
func (w *LoggingWriter) Write(p []byte) (n int, err error) {
	n, err = w.writer.Write(p)
	if err != nil {
		return
	}
	if w.EnableStackTrace {
		var writtenBytes int
		buf := debug.Stack()
		if len(buf) > 0 {
			buf = append(buf, '\n')
		}
		writtenBytes, err = w.writer.Write(buf)
		n = n + writtenBytes
	}
	return
}

func (w *LoggingWriter) fork(writer io.Writer) *LoggingWriter {
	if writer == nil {
		writer = io.Discard
	}

	if loggingWriter, ok := writer.(*LoggingWriter); ok {
		loggingWriter.EnableStackTrace = w.EnableStackTrace
		return loggingWriter
	}

	return &LoggingWriter{
		writer:           writer,
		EnableStackTrace: w.EnableStackTrace,
	}
}
