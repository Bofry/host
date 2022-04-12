package internal

type ComponentService struct {
	components []Runner
}

func NewComponentService() *ComponentService {
	return &ComponentService{}
}

func (m *ComponentService) Start() {
	if m.components != nil {
		for i := 0; i < len(m.components); i++ {
			component := m.components[i]
			component.Start()
		}
	}
}

func (m *ComponentService) Stop() {
	if m.components != nil {
		for i := 0; i < len(m.components); i++ {
			component := m.components[i]
			component.Stop()
		}
	}
}

func (m *ComponentService) RegisterComponent(component Runner) {
	if component != nil {
		m.components = append(m.components, component)
	}
}