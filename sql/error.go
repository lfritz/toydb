package sql

type SyntaxError struct {
	Position int
	Msg      string
}

func (e SyntaxError) Error() string {
	return e.Msg
}
