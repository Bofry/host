package app

import "sync"

var (
	_ SessionStateManager = new(StdSessionStateManager)
)

type StdSessionStateManager struct {
	states map[string]SessionState

	mutex sync.Mutex
}

func NewStdSessionStateManager() *StdSessionStateManager {
	return &StdSessionStateManager{
		states: make(map[string]SessionState),
	}
}

// TryCreate implements SessionStateManager.
func (m *StdSessionStateManager) TryCreate(id string) bool {
	_, ok := m.states[id]
	if !ok {
		m.mutex.Lock()
		defer m.mutex.Unlock()

		_, ok = m.states[id]
		if !ok {
			m.states[id] = NewStdSessionState()
			return true
		}
	}
	return false
}

// Delete implements SessionStateManager.
func (m *StdSessionStateManager) Delete(id string) {
	delete(m.states, id)
}

// Load implements SessionStateManager.
func (m *StdSessionStateManager) Load(id string) SessionState {
	v, ok := m.states[id]
	if !ok {
		return nil
	}
	return v
}

// Update implements SessionStateManager.
func (m *StdSessionStateManager) Update(id string, state SessionState) {
	m.states[id] = state
}
