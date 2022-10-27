package types

func Dec(input string) Value {
	return NewValue(NewDecimal(input))
}

func Txt(input string) Value {
	return NewValue(NewText(input))
}

func Dat(year, month, day int) Value {
	return NewValue(NewDate(year, month, day))
}

func Boo(b bool) Value {
	return NewValue(NewBoolean(b))
}
