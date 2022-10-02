package sql

type Parser[T any] func(tokens *TokenList) (T, *TokenList, error)

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
	token, err := tokens.Peek(TokenTypeString, TokenTypeNumber, TokenTypeIdentifier)
	if err != nil {
		return nil, nil, err
	}
	switch token.Type {
	case TokenTypeString:
		return ParseString(tokens)
	case TokenTypeNumber:
		return ParseNumber(tokens)
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
	decimal, err := ParseDecimal(token.Text)
	if err != nil {
		return nil, nil, SyntaxError{token.From, err.Error()}
	}
	return Number{decimal}, tokens, nil
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
