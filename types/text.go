package types

// Text represents a Unicode string of arbitrary size.
type Text struct {
	value string
}

// NewText returns a new Text instance.
func NewText(value string) Text {
	return Text{value}
}

// Compare compares t to u.
func (t Text) Compare(u Text) Compared {
	switch {
	case t.value < u.value:
		return ComparedLt
	case t.value > u.value:
		return ComparedGt
	}
	return ComparedEq
}
