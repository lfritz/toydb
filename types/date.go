package types

import "fmt"

// Date represents a Gregorian calendar date between year 1 and year 9999.
type Date struct {
	year, month, day int
}

// NewDate returns a new Date instance. It panics if the date is invalid.
func NewDate(year, month, day int) Date {
	result, ok := CheckDate(year, month, day)
	if !ok {
		panic(fmt.Sprintf("invalid date: %04d-%02d-%02d", year, month, day))
	}
	return result
}

// CheckDate returns a new Date instance. If the date is invalid, ok will be false.
func CheckDate(year, month, day int) (date Date, ok bool) {
	if year < 1 || year > 9999 || month < 1 || month > 12 || day < 1 || day > daysInMonth(year, month) {
		return Date{}, false
	}
	return Date{year, month, day}, true
}

func (d Date) Type() Type {
	return TypeDate
}

func daysInMonth(year, month int) int {
	switch month {
	case 1, 3, 5, 7, 8, 10, 12:
		return 31
	case 2:
		if leapYear(year) {
			return 29
		} else {
			return 28
		}
	}
	return 30
}

func leapYear(year int) bool {
	return year%4 == 0 && (year%100 != 0 || year%400 == 0)
}

// Compare compares d to v.
func (d Date) Compare(v BasicValue) Compared {
	e, ok := v.(Date)
	if !ok {
		return ComparedInvalid
	}

	switch {
	case d.year < e.year:
		return ComparedLt
	case d.year > e.year:
		return ComparedGt
	case d.month < e.month:
		return ComparedLt
	case d.month > e.month:
		return ComparedGt
	case d.day < e.day:
		return ComparedLt
	case d.day > e.day:
		return ComparedGt
	}
	return ComparedEq
}

func (d Date) String() string {
	return fmt.Sprintf("%04d-%02d-%02d", d.year, d.month, d.day)
}
