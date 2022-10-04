package sql

import (
	"fmt"
	"strings"
	"unicode"
)

// TODO add "true" and "false" keywords

func Tokenize(input string) ([]Token, error) {
	return NewLexer(input).Run()
}

type LexerState int

const (
	LexerStateStart LexerState = iota
	LexerStateWord
	LexerStatePunctuation
	LexerStateNumber
	LexerStateString
)

type Lexer struct {
	input  []rune
	state  LexerState
	from   int
	next   int
	tokens []Token
	err    error
}

func NewLexer(input string) *Lexer {
	return &Lexer{
		input: []rune(input),
	}
}

func (l *Lexer) Run() ([]Token, error) {
	for {
		r, ok := l.nextRune()
		if !ok {
			return l.tokens, nil
		}

		switch l.state {
		case LexerStateStart:
			switch {
			case isSpace(r):
			case isWordStartCharacter(r):
				l.changeState(LexerStateWord)
			case isDigitOrDot(r):
				l.changeState(LexerStateNumber)
			case isQuote(r):
				l.changeState(LexerStateString)
			case isPunctuation(r):
				l.changeState(LexerStatePunctuation)
			default:
				l.errorf(l.next, "unexpected character: '%c'", r)
			}
		case LexerStateWord:
			switch {
			case isWordCharacter(r):
			case isSpace(r):
				l.tokenForWord()
				l.changeState(LexerStateStart)
			case isPunctuation(r):
				l.tokenForWord()
				l.changeState(LexerStatePunctuation)
			default:
				l.errorf(l.next, "unexpected character: %c", r)
			}
		case LexerStateNumber:
			switch {
			case isDigitOrDot(r):
			case isSpace(r):
				l.tokenForNumber()
				l.changeState(LexerStateStart)
			case isPunctuation(r):
				l.tokenForNumber()
				l.changeState(LexerStatePunctuation)
			default:
				l.errorf(l.next, "unexpected character: %c", r)
			}
		case LexerStateString:
			switch {
			case isQuote(r):
				l.tokenForString()
				l.changeState(LexerStateStart)
			default:
			}
		case LexerStatePunctuation:
			switch {
			case isPunctuation(r):
			case isDigitOrDot(r):
				l.tokenForPunctuation()
				l.changeState(LexerStateNumber)
			case isQuote(r):
				l.tokenForPunctuation()
				l.changeState(LexerStateString)
			case isWordStartCharacter(r):
				l.tokenForPunctuation()
				l.changeState(LexerStateWord)
			case isSpace(r):
				l.tokenForPunctuation()
				l.changeState(LexerStateStart)
			default:
				l.errorf(l.next, "unexpected character: %c", r)
			}
		default:
			panic(fmt.Sprintf("Unexpected state: %v", l.state))
		}

		if l.err != nil {
			return nil, l.err
		}
		l.next++
	}
}

func (l *Lexer) addToken(from, to int, typ TokenType) {
	token := Token{
		Text: string(l.input[from:to]),
		Type: typ,
		From: from,
		To:   to,
	}
	l.tokens = append(l.tokens, token)
}

func (l *Lexer) changeState(s LexerState) {
	l.state = s
	l.from = l.next
}

func (l *Lexer) errorf(position int, format string, a ...any) {
	l.err = SyntaxError{
		Position: position,
		Msg:      fmt.Sprintf(format, a...),
	}
}

func (l *Lexer) tokenForNumber() {
	l.addToken(l.from, l.next, TokenTypeNumber)
}

func (l *Lexer) tokenForString() {
	from := l.from + 1 // skip opening quote
	to := l.next
	l.addToken(from, to, TokenTypeString)
}

func (l *Lexer) tokenForWord() {
	word := string(l.input[l.from:l.next])
	tokenType := TokenTypeIdentifier
	if t, ok := keywordMap[strings.ToLower(word)]; ok {
		tokenType = t
	}
	l.addToken(l.from, l.next, tokenType)
}

func (l *Lexer) tokenForPunctuation() {
	text := string(l.input[l.from:l.next])
	tokenType, ok := punctuationMap[text]
	if !ok {
		l.errorf(l.from, "invalid SQL: %q", text)
		return
	}
	l.addToken(l.from, l.next, tokenType)
}

func (l *Lexer) nextRune() (r rune, ok bool) {
	if l.next > len(l.input) {
		ok = false
		return
	}
	if l.next == len(l.input) {
		// return an extra space at the end to produce the last token
		r = ' '
	} else {
		r = l.input[l.next]
	}
	ok = true
	return
}

func isQuote(r rune) bool {
	return r == '\''
}

func isWordStartCharacter(r rune) bool {
	return unicode.IsLetter(r) || r == '_'
}

func isWordCharacter(r rune) bool {
	return isWordStartCharacter(r) || isDigit(r) || r == '$'
}

func isDigitOrDot(r rune) bool {
	return isDigit(r) || r == '.'
}

func isDigit(r rune) bool {
	return '0' <= r && r <= '9'
}

func isPunctuation(r rune) bool {
	switch r {
	case ',', '.', ';', '=', '!', '<', '>', '(', ')', '*':
		return true
	}
	return false
}

func isSpace(r rune) bool {
	switch r {
	case ' ', '\t', '\r', '\n':
		return true
	}
	return false
}
