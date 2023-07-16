package app

const (
	Nop = appError("Nop")
)

var (
	_ error = appError("")
)

type appError string

func (e appError) Error() string {
	return string(e)
}
