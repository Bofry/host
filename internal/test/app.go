package test

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Bofry/host"
	"github.com/Bofry/trace"
	"go.opentelemetry.io/otel/propagation"
)

var (
	logger *log.Logger = log.New(os.Stdout, "[host-demo] ", log.LstdFlags|log.Lmsgprefix|log.LUTC)
)

var (
	_ host.App                    = new(App)
	_ host.AppStaterConfigurator  = new(App)
	_ host.AppTracingConfigurator = new(App)
)

type (
	App struct {
		Host            *Host
		Config          *Config
		ServiceProvider *ServiceProvider

		Component *MockComponent
	}

	Host struct {
		address  string
		compress bool
	}

	Config struct {
		// server
		ListenAddress  string `arg:"address"`
		EnableCompress bool   `arg:"compress"`

		// redis
		RedisHost     string `env:"*REDIS_HOST"       yaml:"redisHost"`
		RedisPassword string `env:"*REDIS_PASSWORD"   yaml:"redisPassword"`
		RedisDB       int    `env:"REDIS_DB"          yaml:"redisDB"`
		RedisPoolSize int    `env:"REDIS_POOL_SIZE"   yaml:"redisPoolSize"`
		Workspace     string `env:"-"                 yaml:"workspace"`
	}

	ServiceProvider struct {
		RedisClient *MockRedis
	}
)

func (app *App) Init() {
	fmt.Println("App.Init()")

	app.Component = &MockComponent{}
}

func (app *App) OnInit() {
	fmt.Println("App.OnInit()")
}

func (app *App) OnInitComplete() {
	fmt.Println("App.OnInitComplete()")
}

func (app *App) OnStart(ctx context.Context) {
	fmt.Println("App.OnStart()")
}

func (app *App) OnStop(ctx context.Context) {
	fmt.Println("App.OnStop()")
}

func (app *App) ConfigureLogger(logger *log.Logger) {
	fmt.Println("App.ConfigureLogger()")
}

func (app *App) Logger() *log.Logger {
	fmt.Println("App.Logger()")
	return logger
}

func (app *App) ConfigureTracerProvider() {
	fmt.Println("App.ConfigureTracerProvider()")
}

func (app *App) TracerProvider() *trace.SeverityTracerProvider {
	fmt.Println("App.TracerProvider()")
	return nil
}

func (app *App) ConfigureTextMapPropagator() {
	fmt.Println("App.ConfigureTextMapPropagator()")
}

func (app *App) TextMapPropagator() propagation.TextMapPropagator {
	fmt.Println("App.TextMapPropagator()")
	return nil
}

func (provider *ServiceProvider) Init(conf *Config, app *App) {
	provider.RedisClient = &MockRedis{
		Host:     conf.RedisHost,
		Password: conf.RedisPassword,
		DB:       conf.RedisDB,
		PoolSize: conf.RedisPoolSize,
	}
}

func (host *Host) Init(conf *Config) {
	host.address = conf.ListenAddress
	host.compress = conf.EnableCompress
}

func (host *Host) Start(ctx context.Context) {
	logger.Println("[MockApp] Host.Start()")
}

func (host *Host) Stop(ctx context.Context) error {
	logger.Println("[MockApp] Host.Shutdown()")
	return nil
}

func (host *Host) OnError(err error) (disposed bool) {
	return false
}
