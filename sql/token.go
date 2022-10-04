package sql

import (
	"fmt"
	"strings"
)

type TokenType int

const (
	TokenTypeIdentifier TokenType = iota

	// literals
	TokenTypeString
	TokenTypeNumber

	// punctuation
	TokenTypeComma
	TokenTypeDot
	TokenTypeStar
	TokenTypeSemicolon
	TokenTypeOpenParen
	TokenTypeCloseParen
	TokenTypeEq
	TokenTypeNe
	TokenTypeLt
	TokenTypeGt
	TokenTypeLe
	TokenTypeGe

	// keywords
	TokenTypeSelect
	TokenTypeFrom
	TokenTypeWhere
	TokenTypeAnd
	TokenTypeOr
	TokenTypeNot
	TokenTypeIs
	TokenTypeNull
	TokenTypeLeft
	TokenTypeRight
	TokenTypeOuter
	TokenTypeJoin
	TokenTypeOn
	TokenTypeFalse
	TokenTypeTrue
)

var tokenTypeNames = map[TokenType]string{
	TokenTypeIdentifier: "identifier",
	TokenTypeString:     "string",
	TokenTypeNumber:     "number",
	TokenTypeComma:      "comma",
	TokenTypeDot:        "dot",
	TokenTypeStar:       "star",
	TokenTypeSemicolon:  "semicolon",
	TokenTypeOpenParen:  "openparen",
	TokenTypeCloseParen: "closeparen",
	TokenTypeEq:         "eq",
	TokenTypeNe:         "ne",
	TokenTypeLt:         "lt",
	TokenTypeGt:         "gt",
	TokenTypeLe:         "le",
	TokenTypeGe:         "ge",
	TokenTypeSelect:     "select",
	TokenTypeFrom:       "from",
	TokenTypeWhere:      "where",
	TokenTypeAnd:        "and",
	TokenTypeOr:         "or",
	TokenTypeNot:        "not",
	TokenTypeIs:         "is",
	TokenTypeNull:       "null",
	TokenTypeLeft:       "left",
	TokenTypeRight:      "right",
	TokenTypeOuter:      "outer",
	TokenTypeJoin:       "join",
	TokenTypeOn:         "on",
	TokenTypeFalse:      "false",
	TokenTypeTrue:       "true",
}

func (t TokenType) String() string {
	if name, ok := tokenTypeNames[t]; ok {
		return name
	}
	return fmt.Sprintf("unexpected token type: %d", t)
}

var keywordMap = map[string]TokenType{
	"select": TokenTypeSelect,
	"from":   TokenTypeFrom,
	"where":  TokenTypeWhere,
	"and":    TokenTypeAnd,
	"or":     TokenTypeOr,
	"not":    TokenTypeNot,
	"is":     TokenTypeIs,
	"null":   TokenTypeNull,
	"left":   TokenTypeLeft,
	"right":  TokenTypeRight,
	"outer":  TokenTypeOuter,
	"join":   TokenTypeJoin,
	"on":     TokenTypeOn,
	"false":  TokenTypeFalse,
	"true":   TokenTypeTrue,
}

var punctuationMap = map[string]TokenType{
	",":  TokenTypeComma,
	".":  TokenTypeDot,
	"*":  TokenTypeStar,
	";":  TokenTypeSemicolon,
	"(":  TokenTypeOpenParen,
	")":  TokenTypeCloseParen,
	"=":  TokenTypeEq,
	"!=": TokenTypeNe,
	"<>": TokenTypeNe,
	"<":  TokenTypeLt,
	">":  TokenTypeGt,
	"<=": TokenTypeLe,
	">=": TokenTypeGe,
}

type Token struct {
	Type     TokenType
	Text     string
	From, To int
}

func (t Token) String() string {
	if t.Type >= TokenTypeComma {
		return t.Type.String()
	}
	return fmt.Sprintf("(%s %q)", t.Type, t.Text)
}

func PrintTokens(tokens []Token) string {
	values := make([]string, len(tokens))
	for i, t := range tokens {
		values[i] = t.String()
	}
	return strings.Join(values, " ")
}
