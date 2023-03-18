package test

import (
	"context"
	"flag"
	"log"
	"os"
	"testing"
	"time"

	"github.com/Bofry/config"
	"github.com/Bofry/host"
)

func TestStarter(t *testing.T) {
	/* NOTE: panic: CryptAcquireContext: Provider DLL failed to initialize correctly.
	 *
	 * If the following commands applied, the CryptAcquireContext error will be occurred .
	 *  - os.Clearenv()
	 */

	// the following statement like
	// $ export REDIS_HOST=kubernate-redis:26379
	// $ export REDIS_PASSWORD=1234
	// $ export REDIS_POOL_SIZE=128
	initializeEnvironment()

	// the following statement like
	// $ go run app.go --address ":10094" --compress true
	initializeArgs()

	log.Default().SetFlags(log.Default().Flags() | log.LUTC)

	app := App{}
	starter := host.Startup(&app).
		Middlewares().
		ConfigureConfiguration(func(service *config.ConfigurationService) {
			service.
				LoadEnvironmentVariables("").
				LoadYamlFile("config.yaml").
				LoadCommandArguments().
				Output()
		})

	runCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := starter.Start(runCtx); err != nil {
		t.Error(err)
	}

	select {
	case <-runCtx.Done():
		if err := starter.Stop(context.Background()); err != nil {
			t.Error(err)
		}
	}

	// assert app
	{
		if app.Component == nil {
			t.Error("assert 'App.Component':: should not be nil")
		}
	}

	// assert app.Host
	{
		if app.Host == nil {
			t.Error("assert 'App.Host':: should not be nil")
		}
		host := app.Host
		if host.address != ":10094" {
			t.Errorf("assert 'Host.address':: expected '%v', got '%v'", ":10094", host.address)
		}
		if host.compress != true {
			t.Errorf("assert 'Host.compress':: expected '%v', got '%v'", true, host.compress)
		}
	}
	// assert app.Config
	{
		if app.Config == nil {
			t.Error("assert 'App.Config':: should not be nil")
		}
		conf := app.Config
		if conf.ListenAddress != ":10094" {
			t.Errorf("assert 'Config.ListenAddress':: expected '%v', got '%v'", ":10094", conf.ListenAddress)
		}
		if conf.EnableCompress != true {
			t.Errorf("assert 'Config.EnableCompress':: expected '%v', got '%v'", true, conf.EnableCompress)
		}
		if conf.RedisHost != "kubernate-redis:26379" {
			t.Errorf("assert 'Config.RedisHost':: expected '%v', got '%v'", "kubernate-redis:26379", conf.RedisHost)
		}
		if conf.RedisPassword != "1234" {
			t.Errorf("assert 'Config.RedisPassword':: expected '%v', got '%v'", "1234", conf.RedisPassword)
		}
		if conf.RedisDB != 3 {
			t.Errorf("assert 'Config.RedisDB':: expected '%v', got '%v'", 3, conf.RedisDB)
		}
		if conf.RedisPoolSize != 128 {
			t.Errorf("assert 'Config.RedisPoolSize':: expected '%v', got '%v'", 128, conf.RedisPoolSize)
		}
		if conf.Workspace != "demo_test" {
			t.Errorf("assert 'Config.Workspace':: expected '%v', got '%v'", "demo_test", conf.Workspace)
		}
	}
	// assert app.ServiceProvider
	{
		if app.ServiceProvider == nil {
			t.Error("assert 'App.ServiceProvider':: should not be nil")
		}
		provider := app.ServiceProvider
		if provider.RedisClient == nil {
			t.Error("assert 'ServiceProvider.RedisClient':: should not be nil")
		}
		redisClient := provider.RedisClient
		if redisClient.Host != "kubernate-redis:26379" {
			t.Errorf("assert 'RedisClient.Host':: expected '%v', got '%v'", "kubernate-redis:26379", redisClient.Host)
		}
		if redisClient.Password != "1234" {
			t.Errorf("assert 'RedisClient.Password':: expected '%v', got '%v'", "1234", redisClient.Password)
		}
		if redisClient.DB != 3 {
			t.Errorf("assert 'RedisClient.DB':: expected '%v', got '%v'", 3, redisClient.DB)
		}
		if redisClient.PoolSize != 128 {
			t.Errorf("assert 'RedisClient.PoolSize':: expected '%v', got '%v'", 128, redisClient.PoolSize)
		}
	}
}

func initializeEnvironment() {
	os.Setenv("REDIS_HOST", "kubernate-redis:26379")
	os.Setenv("REDIS_PASSWORD", "1234")
	os.Setenv("REDIS_POOL_SIZE", "128")
}

func initializeArgs() {
	os.Args = []string{"example",
		"--address", ":10094",
		"--compress", "true"}

	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
}
