package types

import "fmt"

// Text represents a Unicode string of arbitrary size.
type Text struct {
	value string
}

// NewText returns a new Text instance.
func NewText(value string) Text {
	return Text{value}
}

func (t Text) Type() Type {
	return TypeText
}

// Compare compares t to v.
func (t Text) Compare(v BasicValue) Compared {
	u, ok := v.(Text)
	if !ok {
		return ComparedInvalid
	}
	switch {
	case t.value < u.value:
		return ComparedLt
	case t.value > u.value:
		return ComparedGt
	}
	return ComparedEq
}

func (t Text) String() string {
	return fmt.Sprintf("%q", t.value)
}
