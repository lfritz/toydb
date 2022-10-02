package sql

func ParseExpression(tokens *TokenList) (Expression, *TokenList, error) {
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
