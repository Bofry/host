package internal

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"unsafe"
)

type (
	MockApp struct {
		Config          *MockConfig
		ServiceProvider *MockServiceProvider
		Host            *MyHost
	}

	MyHost MockHost

	MockConfig struct{}

	MockServiceProvider struct{}
)

type MockHost struct {
	v string
}

func (h *MockHost) Start(ctx context.Context) {
	fmt.Printf("MockHost.Start")
}
func (h *MockHost) Stop(ctx context.Context) error { return nil }

func Test(t *testing.T) {
	appCtx := NewAppContext(&MockApp{})
	rvConfig := appCtx.Field(APP_CONFIG_FIELD)
	rvServiceProvider := appCtx.Field(APP_SERVICE_PROVIDER_FIELD)
	rvHost := appCtx.Field(APP_HOST_FIELD)

	var rvHostInterface reflect.Value
	if rvHost.Type().ConvertibleTo(typeOfHost) {
		rvHostInterface = rvHost.Convert(typeOfHost)
		t.Logf("HostInterface1: %+v\n", rvHost.Type().Name())
	} else {
		var typeOfMockHost = reflect.TypeOf(MockHost{})
		rv := reflect.NewAt(typeOfMockHost, unsafe.Pointer(rvHost.Pointer()))
		t.Logf("rv: %#v\n", rv)
		rvHostInterface = rv.Convert(typeOfHost)
		t.Logf("HostInterface2: %#v\n", rvHostInterface)
	}

	t.Logf("Config: %+v\n", rvConfig.Elem().Type().Name())
	t.Logf("ServiceProvider: %+v\n", rvServiceProvider.Elem().Type().Name())
	t.Logf("Host: %+v\n", rvHost.Elem().Type().Name())
	t.Logf("HostInterface: %+v\n", rvHostInterface.Type().Name())
	t.Logf("HostInterface: %+v\n", rvHostInterface.IsNil())

	host, _ := rvHostInterface.Interface().(Host)
	host.Start(context.Background())
}
