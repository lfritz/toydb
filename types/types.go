package types

import "fmt"

// Type defines the types supported by the database.
type Type int

const (
	TypeBoolean Type = iota
	TypeText
	TypeDecimal
	TypeDate
)

func (t Type) String() string {
	switch t {
	case TypeBoolean:
		return "boolean"
	case TypeText:
		return "text"
	case TypeDecimal:
		return "decimal"
	case TypeDate:
		return "date"
	}
	panic(fmt.Sprintf("unexpected Type: %d", t))
}

// The Value interface is implemented by all value types supported by the database.
type Value interface {
	Type() Type
	Compare(v Value) Compared
	String() string
}
