package internal

type HostHelper struct {
	App *AppModule
}

func (h *HostHelper) ExtractHostOnError() HostOnError {
	host := ReflectHelper(h.App.Host()).Value().Interface()
	if host != nil {
		v, ok := host.(HostOnError)
		if ok {
			return v
		}
	}
	return nil
}
