package test

import (
	"context"
	"fmt"
	"log"
)

type (
	MockApp struct {
		Host            *Host
		Config          *Config
		ServiceProvider *ServiceProvider

		Component       *MockComponent
		ComponentRunner *MockComponentRunner
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
		RedisClient *mockRedis
	}
)

func (app *MockApp) Init(conf *Config) {
	fmt.Println("MockApp.Init()")

	app.Component = &MockComponent{}
	app.ComponentRunner = &MockComponentRunner{prefix: "MockComponentRunner"}
}

func (provider *ServiceProvider) Init(conf *Config, app *MockApp) {
	provider.RedisClient = &mockRedis{
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
	log.Println("[MockApp] Host.Start()")
}

func (host *Host) Stop(ctx context.Context) error {
	log.Println("[MockApp] Host.Shutdown()")
	return nil
}
