package app

import (
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/Bofry/host"
	"github.com/Bofry/structproto/reflecting"
)

func Init(v Module, opts ...ModuleBindingOption) *Application {
	var (
		moduleName    string
		buildingOpts  = v.ModuleOptions()
		rvModule      = indirectValue(reflect.ValueOf(v))
		rvApp         reflect.Value
		messageRouter = make(MessageRouter)
		eventRouter   = make(EventRouter)
	)

	if rvModule.Kind() != reflect.Struct {
		panic("Module must be Struct")
	}

	// get module name
	{
		pkgpath := rvModule.Type().PkgPath()
		parts := strings.Split(pkgpath, "/")

		if len(parts) > 0 {
			moduleName = parts[len(parts)-1]
		}
	}

	// get App
	rvApp = rvModule.FieldByName(__MODULE_APP_FIELD)
	if rvApp.IsValid() {
		rvApp = reflecting.AssignZero(rvApp)

		// binding
		for _, opt := range opts {
			err := opt.apply(rvApp)
			if err != nil {
				panic(err)
			}
		}

		// call Init()
		fn := rvApp.MethodByName(host.APP_COMPONENT_INIT_METHOD)
		if fn.IsValid() {
			if fn.Kind() != reflect.Func {
				panic(fmt.Errorf("fail to Init() request handler. cannot find func %s() within type %s\n", host.APP_COMPONENT_INIT_METHOD, rvApp.Type().String()))
			}
			if fn.Type().NumIn() != 0 || fn.Type().NumOut() != 0 {
				panic(fmt.Errorf("fail to Init() request handler. %s.%s() type should be func()\n", rvApp.Type().String(), host.APP_COMPONENT_INIT_METHOD))
			}
			fn.Call([]reflect.Value(nil))
		}

		// binding by fileds
		{
			rvApp := indirectValue(rvApp)

			// binding EventClient
			for i := 0; i < rvApp.Type().NumField(); i++ {
				field := rvApp.Type().Field(i)

				switch field.Name {
				case __APP_EVENT_CLIENT_FIELD:
					if field.Type == typeOfEventClient {
						rvHandler := rvApp.FieldByName(field.Name)
						handler := asEventClient(rvHandler)

						if handler != nil {
							buildingOpts = append(buildingOpts,
								WithEventClient(handler))
						}
					}
				}
			}
		}

		// binding by methods
		{
			// binding DefaultMessageHandler
			{
				rvHandler := rvApp.MethodByName(__APP_DEFAULT_EVENT_HANDLER_METHOD)
				if isEventHandler(rvHandler) {
					handler := asEventHandler(rvHandler)
					if handler != nil {
						buildingOpts = append(buildingOpts,
							WithDefaultEventHandler(handler))
					}
				}
			}

			// binding DefaultEventHandler
			{
				rvHandler := rvApp.MethodByName(__APP_DEFAULT_MESSAGE_HANDLER_METHOD)
				if isMessageHandler(rvHandler) {
					handler := asMessageHandler(rvHandler)
					if handler != nil {
						buildingOpts = append(buildingOpts,
							WithDefaultMessageHandler(handler))
					}
				}
			}
		}

		// binding MessageHandlers & EventHandlers
		for i := 0; i < rvModule.Type().NumField(); i++ {
			field := rvModule.Type().Field(i)

			switch field.Type {
			case typeOfMessageHandler:
				rvHandler := rvApp.MethodByName(field.Name)
				if !isMessageHandler(rvHandler) {
					panic(fmt.Errorf("binding '%s' failed. cannot convert to MessageHandler", field.Name))
				}
				protocol, ok := field.Tag.Lookup(TAG_PROTOCOL)
				if ok && protocol != "-" {
					handler := asMessageHandler(rvHandler)
					messageRouter[protocol] = handler
				}

			case typeOfEventHandler:
				rvHandler := rvApp.MethodByName(field.Name)
				if !isEventHandler(rvHandler) {
					panic(fmt.Errorf("binding '%s' failed. cannot convert to EventHandler", field.Name))
				}
				channel, ok := field.Tag.Lookup(TAG_CHANNEL)
				if ok && channel != "-" {
					handler := asEventHandler(rvHandler)
					eventRouter[channel] = handler
				}

			}
		}
	}

	// add MessageRouter
	if len(messageRouter) > 0 {
		buildingOpts = append(buildingOpts, WithMessageRouter(messageRouter))
	}

	// add EventRouter
	if len(eventRouter) > 0 {
		buildingOpts = append(buildingOpts, WithEventRouter(eventRouter))
	}

	app, err := Build(moduleName, buildingOpts...)
	if err != nil {
		panic(err)
	}
	return app
}

func ModuleOptions(opts ...ApplicationBuildingOption) ModuleOptionCollection {
	return opts
}

func Build(appName string, opts ...ApplicationBuildingOption) (*Application, error) {
	logger := log.New(log.Default().Writer(), "", log.Default().Flags())
	logger.SetPrefix(fmt.Sprintf(__LOGGER_PREFIX_FORMAT, appName))

	app := &Application{
		Name:              appName,
		logger:            logger,
		tracerProvider:    createNoopTracerProvider(),
		textMapPropagator: createNoopTextMapPropagator(),
		eventClient:       NoopEventClient{},
	}
	app.alloc()

	var err error
	for _, builder := range opts {
		err = builder.apply(app)
		if err != nil {
			return nil, err
		}
	}

	app.init()

	return app, nil
}
