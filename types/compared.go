package types

// Compared is the result of a comparison: less than, equal, or greater than.
type Compared int

const (
	ComparedLt Compared = -1
	ComparedEq Compared = 0
	ComparedGt Compared = 1
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
	}
	return "invalid"
}
