package host

import (
	_ "unsafe"
)

//go:linkname Startup github.com/Bofry/host/internal.NewStarter
func Startup(app interface{}) *Starter

//go:linkname RegisterHostService github.com/Bofry/host/internal.RegisterHostService
func RegisterHostService(starter *Starter, service HostService)

//go:linkname StdHostServiceInstance github.com/Bofry/host/internal.StdHostServiceInstance
func StdHostServiceInstance() HostService
