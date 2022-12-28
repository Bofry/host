package internal

import (
	"fmt"
	"time"

	"go.uber.org/fx/fxevent"
)

type DefaultLogger struct {
	Logger *fxevent.ConsoleLogger
}

func (l *DefaultLogger) LogEvent(event fxevent.Event) {
	fmt.Fprintf(l.Logger.W, "%s ", time.Now().UTC().Format("2006-01-02 15:04:05"))
	l.Logger.LogEvent(event)
}
