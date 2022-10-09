package types

// Type defines the types supported by the database.
type Type int

const (
	TypeBoolean Type = iota
	TypeText
	TypeDecimal
	TypeDate
)

// The Value interface is implemented by all value types supported by the database.
type Value interface {
	Type() Type
	Compare(v Value) Compared
}
