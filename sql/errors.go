package sql

// SyntaxError is the error returned by the lexer / parser if the input is not valid.
//
// The position is the index (in characters, not bytes) in the input where an error was detected.
type SyntaxError struct {
	Position int
	Msg      string
}

func (e SyntaxError) Error() string {
	return e.Msg
}
