package sql

import "testing"

var noTokens = []Token{}
var someTokens = []Token{
	Token{Type: TokenTypeSelect, Text: "select"},
	Token{Type: TokenTypeStar, Text: "*"},
	Token{Type: TokenTypeFrom, Text: "from"},
	Token{Type: TokenTypeIdentifier, Text: "foo"},
}

func TestTokenListPeek(t *testing.T) {
	l := TokenList{noTokens}
	_, err := l.Peek()
	if err == nil {
		t.Error("Peek() did not return error for empty list")
	}

	l = TokenList{someTokens}
	got, err := l.Peek()
	if err != nil {
		t.Fatalf("Peek() returned error: %v", err)
	}
	want := Token{Type: TokenTypeSelect, Text: "select"}
	if got != want {
		t.Errorf("Peek() returned %v, want %v", got, want)
	}

	l = TokenList{someTokens}
	_, err = l.Peek(TokenTypeStar)
	if err == nil {
		t.Error("Peek(TokenTypeStar) did not return error")
	}

	l = TokenList{someTokens}
	_, err = l.Peek(TokenTypeStar, TokenTypeIdentifier, TokenTypeWhere)
	if err == nil {
		t.Error("Peek(TokenTypeStar, TokenTypeIdentifier, TokenTypeWhere) did not return error")
	}
}

func TestTokenTypeGet(t *testing.T) {
	l := TokenList{someTokens}
	got, err := l.Get()
	if err != nil {
		t.Fatalf("Get() returned error: %v", err)
	}
	want := Token{Type: TokenTypeSelect, Text: "select"}
	if got != want {
		t.Errorf("Get() returned %v, want %v", got, want)
	}
}

func TestTokenTypeConsume(t *testing.T) {
	l := TokenList{someTokens}
	err := l.Consume()
	if err != nil {
		t.Fatalf("Consume() returned error: %v", err)
	}
}
