package sql

import "errors"

var NotImplemented = errors.New("not yet implemented")

type SyntaxError struct {
	Position int
	Msg      string
}

func (e SyntaxError) Error() string {
	return e.Msg
}
