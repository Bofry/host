package internal

func RegisterHostService(starter *Starter, service HostModule) {
	if service != nil {
		starter.hostModuleBuilder.HostService(service)
	}
}

func StdHostServiceInstance() HostModule {
	return stdHostService
}
