package internal

func RegisterHostService(starter *Starter, host HostModule) {
	if host != nil {
		starter.appServiceBuilder.HostModule(host)
	}
}

func StdHostServiceInstance() HostModule {
	return stdHostModuleSingleton
}
