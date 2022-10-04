package types

import "fmt"

// Decimal represents a decimal number of arbitrary size.
type Decimal struct {
	negative bool
	digits   []uint8
	n        int // number of digits to the left of the dot
}

// DecimalZero returns the number zero.
func DecimalZero() Decimal {
	return Decimal{}
}

// NewDecimal parses the input and panics if it fails.
func NewDecimal(input string) Decimal {
	result, err := ParseDecimal(input)
	if err != nil {
		panic(fmt.Sprintf("invalid decimal number: %s", input))
	}
	return result
}

// ParseDecimal parses the input and returns an error if it's not a valid decimal.
func ParseDecimal(input string) (Decimal, error) {
	runes := []rune(input)
	negative := false
	if len(runes) > 0 && runes[0] == '-' {
		negative = true
		runes = runes[1:]
	}

	if len(runes) == 0 {
		return Decimal{}, fmt.Errorf("not a valid decimal number: %s", input)
	}

	var n int
	var digits []uint8
	gotDot := false
	for i, c := range runes {
		switch {
		case c >= '0' && c <= '9':
			digits = append(digits, uint8(c-'0'))
		case c == '.':
			if gotDot {
				return Decimal{}, fmt.Errorf("not a valid decimal number: %s", input)
			}
			gotDot = true
			n = i
		default:
			return Decimal{}, fmt.Errorf("not a valid decimal number: %s", input)
		}
	}
	if !gotDot {
		n = len(digits)
	}

	decimal := normalize(negative, digits, n)
	return decimal, nil
}

// Compare compares d to e.
func (d Decimal) Compare(e Decimal) Compared {
	gt, lt := ComparedGt, ComparedLt
	switch {
	case d.negative && e.negative:
		gt, lt = lt, gt
	case !d.negative && e.negative:
		return ComparedGt
	case d.negative && !e.negative:
		return ComparedLt
	}

	switch {
	case d.n > e.n:
		return gt
	case d.n < e.n:
		return lt
	}

	minDigits := len(d.digits)
	if len(e.digits) < minDigits {
		minDigits = len(e.digits)
	}
	for i := 0; i < minDigits; i++ {
		switch {
		case d.digits[i] > e.digits[i]:
			return gt
		case d.digits[i] < e.digits[i]:
			return lt
		}
	}

	switch {
	case len(d.digits) > len(e.digits):
		return gt
	case len(d.digits) < len(e.digits):
		return lt
	}

	return ComparedEq
}

func normalize(negative bool, digits []uint8, n int) Decimal {
	for n > 0 && len(digits) > 0 && digits[0] == 0 {
		if n > 0 {
			n--
		}
		digits = digits[1:]
	}
	for len(digits) > n && digits[len(digits)-1] == 0 {
		digits = digits[:len(digits)-1]
	}
	if len(digits) == 0 {
		return DecimalZero()
	}
	return Decimal{negative: negative, digits: digits, n: n}
}
