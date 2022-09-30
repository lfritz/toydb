package sql

func Parse(input string) (Statement, error) {
	tokens, err := Tokenize(input)
	if err != nil {
		return nil, err
	}
	return ParseTokens(&TokenList{input, tokens})
}

func ParseTokens(tokens *TokenList) (Statement, error) {
	err := tokens.Consume(TokenTypeSelect)
	if err != nil {
		return nil, err
	}

	what, tokens, err := parseExpressions(tokens)
	if err != nil {
		return nil, err
	}

	err = tokens.Consume(TokenTypeFrom)
	if err != nil {
		return nil, err
	}

	from, tokens, err := parseFromExpressions(tokens)
	if err != nil {
		return nil, err
	}

	stmt := &SelectStatement{
		What: what,
		From: from,
	}
	return stmt, nil
}

func parseExpressions(tokens *TokenList) ([]Expression, *TokenList, error) {
	var expressions []Expression

	// parse first expression
	expression, tokens, err := parseExpression(tokens)
	if err != nil {
		return nil, nil, err
	}
	expressions = append(expressions, expression)

	// look for comma
	err = tokens.Consume(TokenTypeComma)
	if err != nil {
		return expressions, tokens, nil
	}

	// parse more expressions
	moreExpressions, tokens, err := parseExpressions(tokens)
	if err != nil {
		return nil, nil, err
	}
	expressions = append(expressions, moreExpressions...)

	return expressions, tokens, nil
}

func parseExpression(tokens *TokenList) (Expression, *TokenList, error) {
	token, err := tokens.Get(TokenTypeStar, TokenTypeIdentifier)
	if err != nil {
		return nil, nil, err
	}
	if token.Type == TokenTypeStar {
		return new(Star), tokens, nil
	} else {
		return &Column{Name: token.Text}, tokens, nil
	}
}

func parseFromExpressions(tokens *TokenList) ([]FromExpression, *TokenList, error) {
	token, err := tokens.Get(TokenTypeIdentifier)
	if err != nil {
		return nil, nil, err
	}
	tableName := TableName(token.Text)

	// TODO look for 'left', 'right', or 'join'

	return []FromExpression{tableName}, tokens, nil
}
