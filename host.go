package host

import (
	_ "unsafe"
)

//go:linkname Startup github.com/Bofry/host/internal.NewStarter
func Startup(app interface{}) *Starter

//go:linkname RegisterHostModule github.com/Bofry/host/internal.RegisterHostModule
func RegisterHostModule(starter *Starter, service HostModule)

//go:linkname NopHostServiceInstance github.com/Bofry/host/internal.NopHostServiceInstance
func NopHostServiceInstance() HostModule
