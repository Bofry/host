package app

import (
	"fmt"
	"log"
)

func ModuleOptions(opts ...ApplicationBuildingOption) ModuleOptionCollection {
	return opts
}

func BuildApplication(appName string, opts ...ApplicationBuildingOption) (*Application, error) {
	logger := log.New(log.Default().Writer(), "", log.Default().Flags())
	logger.SetPrefix(fmt.Sprintf(__LOGGER_PREFIX_FORMAT, appName))

	app := &Application{
		Name:              appName,
		logger:            logger,
		tracerProvider:    createNoopTracerProvider(),
		textMapPropagator: createNoopTextMapPropagator(),
		eventSource:       NoopEventSource{},
		messageSource:     NoopMessageSrouce{},
	}
	app.alloc()

	var err error
	for _, builder := range opts {
		err = builder.apply(app)
		if err != nil {
			return nil, err
		}
	}

	return app, nil
}
