package app

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"

	"github.com/Bofry/host"
	"github.com/Bofry/structproto/reflecting"
)

var (
	_RestrictedApplicationBuildingOptions = sliceToMap[string, bool]([]string{
		APP_BUILDING_OPT_DEFAULT_EVENT_HANDLER,
		APP_BUILDING_OPT_DEFAULT_MESSAGE_HANDLER,
		APP_BUILDING_OPT_INVALID_EVENT_HANDLER,
		APP_BUILDING_OPT_INVALID_MESSAGE_HANDLER,
		APP_BUILDING_OPT_ERROR_HANDLER,
	}, func(key string) bool { return true })
)

func Init(v Module, opts ...ModuleBindingOption) *Application {
	var (
		app           *Application
		moduleName    string
		buildingOpts  = v.ModuleOptions()
		rvModule      = indirectValue(reflect.ValueOf(v))
		rvAppModule   reflect.Value
		messageRouter = make(MessageRouter)
		eventRouter   = make(EventRouter)
	)

	if rvModule.Kind() != reflect.Struct {
		panic("Module must be Struct")
	}

	// validate restricted ModuleBindingOption
	for _, opt := range buildingOpts {
		typename := opt.typeName()
		if _, ok := _RestrictedApplicationBuildingOptions[typename]; ok {
			panic(fmt.Sprintf("specified ModuleBindingOption '%s' is restricted", typename))
		}
	}

	// get AppModule
	rvAppModule = rvModule.FieldByName(__MODULE_APP_FIELD)
	if rvAppModule.IsValid() {
		rvAppModule = reflecting.AssignZero(rvAppModule)

		// get module name
		{
			pkgpath := indirectValue(rvAppModule).Type().PkgPath()
			parts := strings.Split(pkgpath, "/")

			if len(parts) > 0 {
				moduleName = parts[len(parts)-1]
			}
		}

		// binding AppModule
		for _, opt := range opts {
			err := opt.apply(rvAppModule, APP)
			if err != nil {
				panic(err)
			}
		}

		// allocate App
		app = allocApplication(moduleName)

		// binding by fields
		{
			// re-define rvAppModule
			var rvAppModule = indirectValue(rvAppModule)

			// binding EventClient
			for i := 0; i < rvAppModule.Type().NumField(); i++ {
				field := rvAppModule.Type().Field(i)

				switch field.Name {
				// case __APP_EVENT_CLIENT_FIELD:
				// 	if field.Type == typeOfEventClient {
				// 		rvHandler := rvAppModule.FieldByName(field.Name)
				// 		handler := asEventClient(rvHandler)

				// 		if handler != nil {
				// 			buildingOpts = append(buildingOpts,
				// 				WithEventClient(handler))
				// 		}
				// 	}
				case __APP_APP_BASE_FIELD:
					if field.Type != typeOfAppBase {
						panic(fmt.Sprintf("field '%s' should be of type %s", field.Name, typeOfAppBase.String()))
					}

					var (
						appBase = &AppBase{application: app}
					)

					rvAppBase := rvAppModule.FieldByName(field.Name)
					rvAppBase = reflecting.AssignZero(rvAppBase)
					rvAppBase.Set(reflect.ValueOf(appBase))
				}
			}
		}

		// call AppModule.Init()
		fn := rvAppModule.MethodByName(host.APP_COMPONENT_INIT_METHOD)
		if fn.IsValid() {
			if fn.Kind() != reflect.Func {
				panic(fmt.Errorf("fail to Init() request handler. cannot find func %s() within type %s\n", host.APP_COMPONENT_INIT_METHOD, rvAppModule.Type().String()))
			}
			if fn.Type().NumIn() != 0 || fn.Type().NumOut() != 0 {
				panic(fmt.Errorf("fail to Init() request handler. %s.%s() type should be func()\n", rvAppModule.Type().String(), host.APP_COMPONENT_INIT_METHOD))
			}
			fn.Call([]reflect.Value(nil))
		}

		// binding by methods
		{
			// binding InvalidEventHandler
			{
				rvHandler := rvAppModule.MethodByName(__APP_INVALID_EVENT_HANDLER_METHOD)
				if isEventHandler(rvHandler) {
					handler := asEventHandler(rvHandler)
					if handler != nil {
						buildingOpts = append(buildingOpts,
							WithDefaultEventHandler(handler))
					}
				}
			}

			// binding InvalidMessageHandler
			{
				rvHandler := rvAppModule.MethodByName(__APP_INVALID_MESSAGE_HANDLER_METHOD)
				if isMessageHandler(rvHandler) {
					handler := asMessageHandler(rvHandler)
					if handler != nil {
						buildingOpts = append(buildingOpts,
							WithDefaultMessageHandler(handler))
					}
				}
			}

			// binding DefaultEventHandler
			{
				rvHandler := rvAppModule.MethodByName(__APP_DEFAULT_EVENT_HANDLER_METHOD)
				if isEventHandler(rvHandler) {
					handler := asEventHandler(rvHandler)
					if handler != nil {
						buildingOpts = append(buildingOpts,
							WithDefaultEventHandler(handler))
					}
				}
			}

			// binding DefaultMessageHandler
			{
				rvHandler := rvAppModule.MethodByName(__APP_DEFAULT_MESSAGE_HANDLER_METHOD)
				if isMessageHandler(rvHandler) {
					handler := asMessageHandler(rvHandler)
					if handler != nil {
						buildingOpts = append(buildingOpts,
							WithDefaultMessageHandler(handler))
					}
				}
			}

			// binding ErrorHandler
			{
				rvHandler := rvAppModule.MethodByName(__APP_ERROR_HANDLER_METHOD)
				if isErrorHandler(rvHandler) {
					handler := asErrorHandler(rvHandler)
					if handler != nil {
						buildingOpts = append(buildingOpts,
							WithErrorHandler(handler))
					}
				}
			}
		}

		// binding MessageHandlers & EventHandlers
		for i := 0; i < rvModule.Type().NumField(); i++ {
			field := rvModule.Type().Field(i)

			switch field.Type {
			case typeOfMessageHandler:
				rvHandler := rvAppModule.MethodByName(field.Name)
				if !isMessageHandler(rvHandler) {
					panic(fmt.Errorf("binding '%s' failed. cannot convert to MessageHandler", field.Name))
				}
				protocol, ok := field.Tag.Lookup(TAG_PROTOCOL)
				if ok && protocol != "-" {
					handler := asMessageHandler(rvHandler)
					_, ok := messageRouter[protocol]
					if ok {
						panic(fmt.Errorf("find duplicate protocol '%s' on field '%s'", protocol, field.Name))
					}
					messageRouter[protocol] = handler
				}
			case typeOfEventHandler:
				rvHandler := rvAppModule.MethodByName(field.Name)
				if !isEventHandler(rvHandler) {
					panic(fmt.Errorf("binding '%s' failed. cannot convert to EventHandler", field.Name))
				}
				channel, ok := field.Tag.Lookup(TAG_CHANNEL)
				if ok && channel != "-" {
					optExpandEnv := field.Tag.Get(TAG_OPT_EXPAND_ENV)
					if optExpandEnv != OPT_OFF || len(optExpandEnv) == 0 || optExpandEnv == OPT_ON {
						channel = os.ExpandEnv(channel)
					}

					handler := asEventHandler(rvHandler)
					_, ok := eventRouter[channel]
					if ok {
						panic(fmt.Errorf("find duplicate channel '%s' on field '%s'", channel, field.Name))
					}
					eventRouter[channel] = handler
				}
			}
		}
	}

	// binding module options
	for _, opt := range opts {
		err := opt.apply(reflect.ValueOf(&buildingOpts), MODULE_OPTIONS)
		if err != nil {
			panic(err)
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

	// initialize Application
	{
		var err error
		for _, builder := range buildingOpts {
			err = builder.apply(app)
			if err != nil {
				panic(err)
			}
		}
		app.init()
	}
	return app
}

func ModuleOptions(opts ...ApplicationBuildingOption) ModuleOptionCollection {
	return opts
}

func Build(appName string, opts ...ApplicationBuildingOption) (*Application, error) {
	app := allocApplication(appName)

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

func allocApplication(appName string) *Application {
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

	return app
}
