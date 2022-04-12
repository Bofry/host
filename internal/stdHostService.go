package internal

import "reflect"

var _ HostService = new(StdHostService)

type StdHostService struct{}

func (s *StdHostService) Init(host Host, app *AppContext)         {}
func (s *StdHostService) InitComplete(host Host, app *AppContext) {}
func (s *StdHostService) GetHostType() reflect.Type {
	return typeOfHost
}
