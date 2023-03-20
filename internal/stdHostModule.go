package internal

import (
	"io"
	"reflect"
)

var _ HostModule = new(StdHostModule)

type StdHostModule struct{}

func (s *StdHostModule) Init(host Host, app *AppModule)         {}
func (s *StdHostModule) InitComplete(host Host, app *AppModule) {}
func (s *StdHostModule) DescribeHostType() reflect.Type {
	return typeOfHost
}
func (s *StdHostModule) ConfigureLogger(logflags int, w io.Writer) {}
