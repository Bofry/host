package app

type MessageRouter map[string]MessageHandler

func (r MessageRouter) Add(code string, handler MessageHandler) {
	r[code] = handler
}

func (r MessageRouter) Remove(code string) {
	delete(r, code)
}

func (r MessageRouter) Get(code string) MessageHandler {
	if r == nil {
		return nil
	}

	if v, ok := r[code]; ok {
		return v
	}
	return nil
}

func (r MessageRouter) Has(code string) bool {
	if r == nil {
		return false
	}

	if _, ok := r[code]; ok {
		return true
	}
	return false
}
