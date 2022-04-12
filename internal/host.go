package internal

func RegisterHostService(starter *Starter, service HostService) {
	if service != nil {
		starter.hostModuleBuilder.HostService(service)
	}
}

func StdHostServiceInstance() HostService {
	return stdHostService
}
