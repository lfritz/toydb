package sql

import (
	"fmt"
	"strings"
)

type TokenList struct {
	input  string
	tokens []Token
}

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
