package types

// Compared is the result of a comparison: less than, equal, or greater than.
type Compared int

const (
	ComparedLt Compared = iota
	ComparedEq
	ComparedGt
	ComparedNull
	ComparedInvalid
)

func (c Compared) String() string {
	switch c {
	case ComparedLt:
		return "lt"
	case ComparedEq:
		return "eq"
	case ComparedGt:
		return "gt"
	case ComparedNull:
		return "null"
	}
	return "invalid"
}
