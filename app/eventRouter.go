package app

type EventRouter map[string]EventHandler

func (r EventRouter) Add(channel string, handler EventHandler) {
	r[channel] = handler
}

func (r EventRouter) Remove(channel string) {
	delete(r, channel)
}

func (r EventRouter) Get(channel string) EventHandler {
	if r == nil {
		return nil
	}

	if v, ok := r[channel]; ok {
		return v
	}
	return nil
}

func (r EventRouter) Has(channel string) bool {
	if r == nil {
		return false
	}

	if _, ok := r[channel]; ok {
		return true
	}
	return false
}
