package internal

func RegisterHostModule(starter *Starter, host HostModule) {
	if host != nil {
		starter.appServiceBuilder.HostModule(host)
	}
}

func NopHostServiceInstance() HostModule {
	return nopHostModuleSingleton
}
