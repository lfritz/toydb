package sql

import "testing"

var noTokens = &TokenList{
	input:  "",
	tokens: []Token{},
}

var someTokens = &TokenList{
	input: "select * from foo",
	tokens: []Token{
		Token{Type: TokenTypeSelect, Text: "select"},
		Token{Type: TokenTypeStar, Text: "*"},
		Token{Type: TokenTypeFrom, Text: "from"},
		Token{Type: TokenTypeIdentifier, Text: "foo"},
	},
}

func TestTokenListLen(t *testing.T) {
	cases := []struct {
		tokens *TokenList
		want   int
	}{
		{noTokens, 0},
		{someTokens, 4},
	}
	for _, c := range cases {
		got := c.tokens.Len()
		if got != c.want {
			t.Errorf("Len returned %d, want %d", got, c.want)
		}
	}
}

func TestTokenListPeek(t *testing.T) {
	l := noTokens
	_, err := l.Peek()
	if err == nil {
		t.Error("Peek() did not return error for empty list")
	}

	l = someTokens
	got, err := l.Peek()
	if err != nil {
		t.Fatalf("Peek() returned error: %v", err)
	}
	want := Token{Type: TokenTypeSelect, Text: "select"}
	if got != want {
		t.Errorf("Peek() returned %v, want %v", got, want)
	}

	l = someTokens
	_, err = l.Peek(TokenTypeStar)
	if err == nil {
		t.Error("Peek(TokenTypeStar) did not return error")
	}

	l = someTokens
	_, err = l.Peek(TokenTypeStar, TokenTypeIdentifier, TokenTypeWhere)
	if err == nil {
		t.Error("Peek(TokenTypeStar, TokenTypeIdentifier, TokenTypeWhere) did not return error")
	}
}

func TestTokenListGet(t *testing.T) {
	l := someTokens
	got, err := l.Get()
	if err != nil {
		t.Fatalf("Get() returned error: %v", err)
	}
	want := Token{Type: TokenTypeSelect, Text: "select"}
	if got != want {
		t.Errorf("Get() returned %v, want %v", got, want)
	}
}

func TestTokenListConsume(t *testing.T) {
	l := someTokens
	err := l.Consume()
	if err != nil {
		t.Fatalf("Consume() returned error: %v", err)
	}
}

func TestJoinWithOr(t *testing.T) {
	cases := []struct {
		items []TokenType
		want  string
	}{
		{[]TokenType{TokenTypeAnd, TokenTypeOr, TokenTypeNot}, "and, or or not"},
		{[]TokenType{TokenTypeOr, TokenTypeNot}, "or or not"},
		{[]TokenType{TokenTypeNot}, "not"},
	}
	for _, c := range cases {
		got := joinWithOr(c.items)
		if got != c.want {
			t.Errorf("joinWithOr(%#v) == %q, want %q", c.items, got, c.want)
		}
	}
}
