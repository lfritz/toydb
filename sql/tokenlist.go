package sql

import (
	"fmt"
	"strings"
)

// TokenList is the input to the parser: a list of tokens, with the original input text.
type TokenList struct {
	input  string
	tokens []Token
}

// Len returns the number of tokens.
func (l *TokenList) Len() int {
	return len(l.tokens)
}

// Peek returns the next token without modifying the list.
//
// If there are no tokens left of the next token doesn't have one of the expected types, it returns
// an error. If no expected types are given, then any type is accepted.
func (l *TokenList) Peek(expected ...TokenType) (Token, error) {
	if err := l.checkEnd(); err != nil {
		return Token{}, err
	}
	first := l.tokens[0]
	err := l.checkType(expected)
	if err != nil {
		return Token{}, err
	}
	return first, nil
}

// Get removes the next token from the list and returns it. The arguments are used in the same way
// as for Peek.
func (l *TokenList) Get(expected ...TokenType) (Token, error) {
	if err := l.checkEnd(); err != nil {
		return Token{}, err
	}
	first, remaining := l.tokens[0], l.tokens[1:]
	err := l.checkType(expected)
	if err != nil {
		return Token{}, err
	}
	l.tokens = remaining
	return first, nil
}

// Consume removes the next token from the list. The arguments are used in the same way as for Peek.
func (l *TokenList) Consume(expected ...TokenType) error {
	if err := l.checkEnd(); err != nil {
		return err
	}
	err := l.checkType(expected)
	if err != nil {
		return err
	}
	l.tokens = l.tokens[1:]
	return nil
}

// End returns an error if the list is not empty.
func (l *TokenList) ExpectEnd() error {
	if len(l.tokens) > 0 {
		token := l.tokens[0]
		return SyntaxError{
			Position: len([]rune(l.input)),
			Msg:      fmt.Sprintf("got %s, expected end of statement", token.Text),
		}
	}
	return nil
}

func (l *TokenList) checkEnd() error {
	if len(l.tokens) == 0 {
		return SyntaxError{
			Position: len([]rune(l.input)),
			Msg:      "unexpected end of input",
		}
	}
	return nil
}

func (l *TokenList) checkType(expected []TokenType) error {
	token := l.tokens[0]
	if len(expected) == 0 {
		return nil
	}
	for _, e := range expected {
		if token.Type == e {
			return nil
		}
	}
	if len(expected) == 1 {
		return SyntaxError{
			Position: token.From,
			Msg:      fmt.Sprintf("got %q, expected %s", token.Text, expected[0]),
		}
	}
	return SyntaxError{
		Position: token.From,
		Msg:      fmt.Sprintf("got %q, expected %s", token.Text, joinWithOr(expected)),
	}
}

func joinWithOr(items []TokenType) string {
	builder := new(strings.Builder)
	last := len(items) - 1
	for i, item := range items {
		switch {
		case i == 0:
			fmt.Fprintf(builder, "%v", item)
		case i == last:
			fmt.Fprintf(builder, " or %v", item)
		default:
			fmt.Fprintf(builder, ", %v", item)
		}
	}
	return builder.String()
}
