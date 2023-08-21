package test

import (
	"context"
	"flag"
	"log"
	"os"
	"reflect"
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

	runCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
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
		expectedConfig := Config{
			ListenAddress:  ":10094",
			EnableCompress: true,
			RedisHost:      "kubernate-redis:26379",
			RedisPassword:  "1234",
			RedisDB:        3,
			RedisPoolSize:  128,
			Workspace:      "demo_test",
		}
		if !reflect.DeepEqual(expectedConfig, *app.Config) {
			t.Errorf("assert 'Config':: expected '%#+v', got '%#+v'", expectedConfig, app.Config)
		}
	}
	// assert app.ServiceProvider
	{
		expectedServiceProvider := ServiceProvider{
			RedisClient: &MockRedis{
				Host:     "kubernate-redis:26379",
				Password: "1234",
				DB:       3,
				PoolSize: 128,
			},
		}
		if !reflect.DeepEqual(expectedServiceProvider, *app.ServiceProvider) {
			t.Errorf("assert 'ServiceProvider':: expected '%#+v', got '%#+v'", expectedServiceProvider, app.ServiceProvider)
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
