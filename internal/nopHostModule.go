package internal

import (
	"io"
	"reflect"
)

var _ HostModule = NopHostModule{}

type NopHostModule struct{}

func (s NopHostModule) Init(host Host, app *AppModule)         {}
func (s NopHostModule) InitComplete(host Host, app *AppModule) {}
func (s NopHostModule) DescribeHostType() reflect.Type {
	return typeOfHost
}
func (s NopHostModule) ConfigureLogger(logflags int, w io.Writer) {}
func (s NopHostModule) OnError(err error)                         {}
