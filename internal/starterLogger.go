package internal

import (
	"log"
	"time"

	"go.uber.org/fx/fxevent"
)

type StarterLogger struct {
	Flags  int
	Logger *fxevent.ConsoleLogger
}

func (l *StarterLogger) LogEvent(event fxevent.Event) {
	var buffer []byte
	t := time.Now()

	if l.Flags&(log.Ldate|log.Ltime|log.Lmicroseconds) != 0 {
		if l.Flags&log.LUTC != 0 {
			t = t.UTC()
		}
		var buf = &buffer
		if l.Flags&log.Ldate != 0 {
			year, month, day := t.Date()
			itoa(buf, year, 4)
			*buf = append(*buf, '/')
			itoa(buf, int(month), 2)
			*buf = append(*buf, '/')
			itoa(buf, day, 2)
			*buf = append(*buf, ' ')
		}
		if l.Flags&(log.Ltime|log.Lmicroseconds) != 0 {
			hour, min, sec := t.Clock()
			itoa(buf, hour, 2)
			*buf = append(*buf, ':')
			itoa(buf, min, 2)
			*buf = append(*buf, ':')
			itoa(buf, sec, 2)
			*buf = append(*buf, ' ')
		}
	}
	l.Logger.W.Write(buffer)
	l.Logger.LogEvent(event)
}
