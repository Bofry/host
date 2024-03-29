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

func (h *MyHost) OnError(err error) (disposed bool) {
	return false
}

func TestAppModule(t *testing.T) {
	var (
		app               = NewAppModule(&MockApp{})
		rvConfig          = app.Config()
		rvServiceProvider = app.ServiceProvider()
		rvHost            = app.Host()
	)

	// TODO: add test assertion
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

	hostHelper := &AppHostHelper{
		App: app,
	}
	handler := hostHelper.OnErrorEventHandler()
	t.Logf("OnErrorEventHandler: %+v\n", handler)
	host.Start(context.Background())
}
