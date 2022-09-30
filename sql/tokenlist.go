package sql

import (
	"errors"
	"fmt"
)

type TokenList struct {
	tokens []Token
}

func (l *TokenList) Peek(expected ...TokenType) (Token, error) {
	if len(l.tokens) == 0 {
		return Token{}, errors.New("unexpected end of input")
	}
	first := l.tokens[0]
	err := checkType(first, expected)
	if err != nil {
		return Token{}, err
	}
	return first, nil
}

func (l *TokenList) Get(expected ...TokenType) (Token, error) {
	if len(l.tokens) == 0 {
		return Token{}, errors.New("unexpected end of input")
	}
	first, remaining := l.tokens[0], l.tokens[1:]
	err := checkType(first, expected)
	if err != nil {
		return Token{}, err
	}
	l.tokens = remaining
	return first, nil
}

func (l *TokenList) Consume(expected ...TokenType) error {
	if len(l.tokens) == 0 {
		return errors.New("unexpected end of input")
	}
	err := checkType(l.tokens[0], expected)
	if err != nil {
		return err
	}
	l.tokens = l.tokens[1:]
	return nil
}

func checkType(token Token, expected []TokenType) error {
	if len(expected) == 0 {
		return nil
	}
	got := token.Type
	for _, e := range expected {
		if got == e {
			return nil
		}
	}
	if len(expected) == 1 {
		return fmt.Errorf("got %s, expected %s", got, expected[0])
	}
	return fmt.Errorf("got %s, expected one of: %v", got, expected)
}
