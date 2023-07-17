package app

var (
	_ SessionState = new(StdSessionState)
)

type StdSessionState struct {
	values map[interface{}]interface{}
}

func NewStdSessionState() *StdSessionState {
	return &StdSessionState{
		values: make(map[interface{}]interface{}),
	}
}

// CanVisit implements SessionState.
func (state *StdSessionState) CanVisit() bool {
	return true
}

// Lock implements SessionState.
func (state *StdSessionState) Lock() {}

// SetValue implements SessionState.
func (state *StdSessionState) SetValue(k interface{}, v interface{}) {
	state.values[k] = v
}

// Unlock implements SessionState.
func (state *StdSessionState) Unlock() {}

// Value implements SessionState.
func (state *StdSessionState) Value(k interface{}) interface{} {
	return state.values[k]
}

// Visit implements SessionState.
func (state *StdSessionState) Visit(fn func(k interface{}, v interface{})) {
	for k, v := range state.values {
		fn(k, v)
	}
}
