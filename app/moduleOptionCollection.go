package app

type ModuleOptionCollection []ApplicationBuildingOption

func (c ModuleOptionCollection) ModuleOptions() []ApplicationBuildingOption {
	return c
}
