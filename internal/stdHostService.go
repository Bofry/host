package internal

import (
	"io"
	"reflect"
)

var _ HostService = new(StdHostService)

type StdHostService struct{}

func (s *StdHostService) Init(host Host, app *AppModule)         {}
func (s *StdHostService) InitComplete(host Host, app *AppModule) {}
func (s *StdHostService) DescribeHostType() reflect.Type {
	return typeOfHost
}
func (s *StdHostService) ConfigureLogger(logflags int, w io.Writer) {}
