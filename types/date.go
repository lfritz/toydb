package types

// Date represents a Gregorian calendar date between year 1 and year 9999.
type Date struct {
	year, month, day int
}

// NewDate returns a new Date instance. If the date is invalid, ok will be false.
func NewDate(year, month, day int) (date Date, ok bool) {
	if year < 1 || year > 9999 || month < 1 || month > 12 || day < 1 || day > daysInMonth(year, month) {
		return Date{}, false
	}
	return Date{year, month, day}, true
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

// Compare compares d to e.
func (d Date) Compare(e Date) Compared {
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
