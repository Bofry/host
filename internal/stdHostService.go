package internal

import (
	"log"
	"reflect"
)

var _ HostService = new(StdHostService)

type StdHostService struct{}

func (s *StdHostService) Init(host Host, app *AppContext)         {}
func (s *StdHostService) InitComplete(host Host, app *AppContext) {}
func (s *StdHostService) DescribeHostType() reflect.Type {
	return typeOfHost
}
func (s *StdHostService) ConfigureLogger(logger *log.Logger) {}
