package types

// The BasicValue interface is implemented by the basic types of the database.
type BasicValue interface {
	Type() Type
	Compare(v BasicValue) Compared
	String() string
}

// A Value is a value in the database, which may be null.
type Value struct {
	t    Type
	null bool
	v    BasicValue
}

func NewNull(t Type) Value {
	return Value{t: t, null: true}
}

func NewValue(v BasicValue) Value {
	return Value{t: v.Type(), v: v}
}

func (v Value) Type() Type {
	return v.t
}

func (v Value) Value() BasicValue {
	return v.v
}

func (v Value) Null() bool {
	return v.null
}

func (v Value) IsTrue() bool {
	if v.null {
		return false
	}
	b, ok := v.v.(Boolean)
	if !ok {
		return false
	}
	return b.value
}

func (v Value) Compare(w Value) Compared {
	if v.t != w.t {
		return ComparedInvalid
	}
	if v.null || w.null {
		return ComparedNull
	}
	return v.v.Compare(w.v)
}

func (v Value) String() string {
	if v.null {
		return "null"
	}
	return v.v.String()
}
