package sql

import (
	"regexp"

	"github.com/lfritz/toydb/types"
)

// Parse parses a select statement with optional semicolon at the end.
func Parse(input string) (*SelectStatement, error) {
	ts, err := Tokenize(input)
	if err != nil {
		return nil, err
	}
	tokens := &TokenList{input, ts}

	statement, tokens, err := ParseSelectStatement(tokens)
	if err != nil {
		return nil, err
	}
	_ = tokens.Consume(TokenTypeSemicolon)

	if err = tokens.ExpectEnd(); err != nil {
		return nil, err
	}

	return statement, nil
}

type Parser[T any] func(tokens *TokenList) (T, *TokenList, error)

func ParseSelectStatement(tokens *TokenList) (*SelectStatement, *TokenList, error) {
	err := tokens.Consume(TokenTypeSelect)
	if err != nil {
		return nil, nil, err
	}
	result := new(SelectStatement)

	result.What, tokens, err = ParseSelectList(tokens)
	if err != nil {
		return nil, nil, err
	}

	err = tokens.Consume(TokenTypeFrom)
	if err != nil {
		return nil, nil, err
	}

	result.From, tokens, err = ParseTableReference(tokens)
	if err != nil {
		return nil, nil, err
	}

	err = tokens.Consume(TokenTypeWhere)
	if err == nil {
		result.Where, tokens, err = ParseExpression(tokens)
		if err != nil {
			return nil, nil, err
		}
	}

	return result, tokens, nil
}

func ParseTableReference(tokens *TokenList) (TableReference, *TokenList, error) {
	left, tokens, err := ParseTableName(tokens)
	if err != nil {
		return nil, nil, err
	}

	token, err := tokens.Get(TokenTypeLeft, TokenTypeRight, TokenTypeJoin)
	if err != nil {
		// this is not a join, it's just a single table name
		return left, tokens, nil
	}

	join := &Join{
		Left: left,
	}

	switch token.Type {
	case TokenTypeLeft:
		join.Type = JoinTypeLeftOuter
	case TokenTypeRight:
		join.Type = JoinTypeRightOuter
	}

	switch token.Type {
	case TokenTypeLeft, TokenTypeRight:
		_ = tokens.Consume(TokenTypeOuter)
		err := tokens.Consume(TokenTypeJoin)
		if err != nil {
			return nil, nil, err
		}
	}

	right, tokens, err := ParseTableName(tokens)
	if err != nil {
		return nil, nil, err
	}
	join.Right = right

	err = tokens.Consume(TokenTypeOn)
	if err != nil {
		return nil, nil, err
	}

	expression, tokens, err := ParseExpression(tokens)
	if err != nil {
		return nil, nil, err
	}
	join.Condition = expression

	return join, tokens, nil
}

func ParseTableName(tokens *TokenList) (TableName, *TokenList, error) {
	token, err := tokens.Get(TokenTypeIdentifier)
	if err != nil {
		return TableName{}, nil, err
	}
	return TableName{token.Text}, tokens, nil
}

func ParseSelectList(tokens *TokenList) (SelectList, *TokenList, error) {
	err := tokens.Consume(TokenTypeStar)
	if err == nil {
		return Star{}, tokens, nil
	}

	expressions, tokens, err := ParseExpressionList(tokens)
	if err != nil {
		return nil, nil, err
	}
	return ExpressionList{Expressions: expressions}, tokens, nil
}

func ParseExpressionList(tokens *TokenList) ([]Expression, *TokenList, error) {
	var result []Expression

	first := true
	for {
		e, tokens, err := ParseExpression(tokens)
		if err != nil {
			if first {
				// empty expression list is allowed
				break
			}
			return nil, nil, err
		}
		first = false

		result = append(result, e)

		err = tokens.Consume(TokenTypeComma)
		if err != nil {
			break
		}
	}

	return result, tokens, nil
}

func ParseExpression(tokens *TokenList) (Expression, *TokenList, error) {
	left, tokens, err := ParseValue(tokens)
	if err != nil {
		return nil, nil, err
	}

	token, err := tokens.Get(TokenTypeEq, TokenTypeNe, TokenTypeLt, TokenTypeGt, TokenTypeLe, TokenTypeGe)
	if err != nil {
		return left, tokens, nil
	}
	op := tokenToOperator[token.Type]

	right, tokens, err := ParseValue(tokens)
	if err != nil {
		return nil, nil, err
	}

	result := &BinaryOperation{
		Left:     left,
		Operator: op,
		Right:    right,
	}
	return result, tokens, nil
}

var tokenToOperator = map[TokenType]BinaryOperator{
	TokenTypeEq: BinaryOperatorEq,
	TokenTypeNe: BinaryOperatorNe,
	TokenTypeLt: BinaryOperatorLt,
	TokenTypeGt: BinaryOperatorGt,
	TokenTypeLe: BinaryOperatorLe,
	TokenTypeGe: BinaryOperatorGe,
}

func ParseValue(tokens *TokenList) (Expression, *TokenList, error) {
	token, err := tokens.Peek(
		TokenTypeString,
		TokenTypeNumber,
		TokenTypeIdentifier,
		TokenTypeFalse,
		TokenTypeTrue,
		TokenTypeDate,
	)
	if err != nil {
		return nil, nil, err
	}
	switch token.Type {
	case TokenTypeString:
		return ParseString(tokens)
	case TokenTypeNumber:
		return ParseNumber(tokens)
	case TokenTypeFalse, TokenTypeTrue:
		tokens.Consume()
		return Boolean{token.Type == TokenTypeTrue}, tokens, nil
	case TokenTypeDate:
		tokens.Consume()
		return ParseDate(tokens)
	default:
		return ParseColumnReference(tokens)
	}
}

func ParseString(tokens *TokenList) (Expression, *TokenList, error) {
	token, err := tokens.Get(TokenTypeString)
	if err != nil {
		return nil, nil, err
	}
	return String{token.Text}, tokens, nil
}

func ParseNumber(tokens *TokenList) (Expression, *TokenList, error) {
	token, err := tokens.Get(TokenTypeNumber)
	if err != nil {
		return nil, nil, err
	}
	decimal, err := types.ParseDecimal(token.Text)
	if err != nil {
		return nil, nil, SyntaxError{token.From, err.Error()}
	}
	return Number{decimal}, tokens, nil
}

var dateRe = regexp.MustCompile(`^(\d\d)-(\d\d)-(\d\d\d\d)$`)

func ParseDate(tokens *TokenList) (Expression, *TokenList, error) {
	token, err := tokens.Get(TokenTypeString)
	if err != nil {
		return nil, nil, err
	}
	date, err := types.ParseDate(token.Text)
	if err != nil {
		return nil, nil, SyntaxError{token.From, err.Error()}
	}
	return Date{date}, tokens, nil
}

func ParseColumnReference(tokens *TokenList) (Expression, *TokenList, error) {
	first, err := tokens.Get(TokenTypeIdentifier)
	if err != nil {
		return nil, nil, err
	}

	err = tokens.Consume(TokenTypeDot)
	if err != nil {
		return ColumnReference{Name: first.Text}, tokens, nil
	}

	second, err := tokens.Get(TokenTypeIdentifier)
	if err != nil {
		return nil, nil, err
	}

	result := ColumnReference{
		Relation: first.Text,
		Name:     second.Text,
	}
	return result, tokens, nil
}
