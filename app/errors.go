package app

const (
	Nop = appError("Nop")
)

var (
	_ error = appError("")
	_ error = MessageError(nil)
)

type appError string

func (e appError) Error() string {
	return string(e)
}

func (e appError) String() string {
	return e.Error()
}

type MessageError error
