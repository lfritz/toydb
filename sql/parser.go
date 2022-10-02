package sql

func ParseExpression(tokens *TokenList) (Expression, *TokenList, error) {
	token, err := tokens.Get(TokenTypeString)
	if err != nil {
		return nil, nil, err
	}
	return String{token.Text}, tokens, nil
}
