package test

import (
	"context"
	"fmt"
	"log"
	"os"
)

var (
	logger *log.Logger = log.New(os.Stdout, "", log.LstdFlags|log.LUTC)
)

type (
	MockApp struct {
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

func (app *MockApp) Init(conf *Config) {
	fmt.Println("MockApp.Init()")

	app.Component = &MockComponent{}
}

func (provider *ServiceProvider) Init(conf *Config, app *MockApp) {
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
