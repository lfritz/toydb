package sql

type Parser[T any] func(tokens *TokenList) (T, *TokenList, error)

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
	token, err := tokens.Get(TokenTypeString, TokenTypeNumber)
	if err != nil {
		return nil, nil, err
	}

	var exp Expression
	switch token.Type {
	case TokenTypeString:
		exp = String{token.Text}
	case TokenTypeNumber:
		decimal, err := ParseDecimal(token.Text)
		if err != nil {
			return nil, nil, SyntaxError{token.From, err.Error()}
		}
		exp = Number{decimal}
	}

	return exp, tokens, nil
}
