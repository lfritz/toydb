package types

// Boolean represents a boolean value.
type Boolean struct {
	value bool
}

// NewBoolean returns a new Boolean instance.
func NewBoolean(value bool) Boolean {
	return Boolean{value}
}

// Compare compares b to c, with true > false.
func (b Boolean) Compare(c Boolean) Compared {
	switch {
	case b.value == c.value:
		return ComparedEq
	case b.value == false:
		return ComparedLt
	}
	return ComparedGt
}
