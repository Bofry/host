package internal

import "log"

type ComponentService struct {
	components []Runable
	logger     *log.Logger
}

func NewComponentService(logger *log.Logger) *ComponentService {
	return &ComponentService{
		logger: logger,
	}
}

func (m *ComponentService) Start() {
	if m.components != nil {
		for i := 0; i < len(m.components); i++ {
			component := m.components[i]
			m.logger.Printf("STARTING Component %T", component)
			component.Runner().Start()
		}
	}
}

func (m *ComponentService) Stop() {
	if m.components != nil {
		for i := 0; i < len(m.components); i++ {
			component := m.components[i]
			component.Runner().Stop()
			m.logger.Printf("STOPPED Component %T", component)
		}
	}
}

func (m *ComponentService) RegisterComponent(component Runable) {
	if component != nil {
		m.components = append(m.components, component)
	}
}
