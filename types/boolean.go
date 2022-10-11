package types

// Boolean represents a boolean value.
type Boolean struct {
	value bool
}

// NewBoolean returns a new Boolean instance.
func NewBoolean(value bool) Boolean {
	return Boolean{value}
}

func (b Boolean) Type() Type {
	return TypeBoolean
}

func (b Boolean) Bool() bool {
	return b.value
}

// Compare compares b to v, with true > false.
func (b Boolean) Compare(v Value) Compared {
	c, ok := v.(Boolean)
	if !ok {
		return ComparedInvalid
	}
	switch {
	case b.value == c.value:
		return ComparedEq
	case b.value == false:
		return ComparedLt
	}
	return ComparedGt
}
