package internal

type AppHostHelper struct {
	App *AppModule
}

func (h *AppHostHelper) OnErrorEventHandler() HostOnErrorEventHandler {
	host := ReflectHelper(h.App.Host()).Value().Interface()
	if host != nil {
		v, ok := host.(HostOnErrorEventHandler)
		if ok {
			return v
		}
	}
	return nil
}
