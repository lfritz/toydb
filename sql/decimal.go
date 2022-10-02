package sql

import (
	"fmt"
	"strings"
)

// Decimal is a quick-and-dirty implementation of a decimal numberic type.
type Decimal struct {
	Value, Digits int
}

func ParseDecimal(input string) (Decimal, error) {
	if input == "" {
		return Decimal{}, fmt.Errorf("not a valid decimal number: %s", input)
	}
	before, after, _ := strings.Cut(input, ".")

	x, _, negative, ok := parseInt(before)
	if !ok {
		return Decimal{}, fmt.Errorf("not a valid decimal number: %s", input)
	}

	after = strings.TrimRight(after, "0")
	y, digits, negativeFractionalPart, ok := parseInt(after)
	if negativeFractionalPart || !ok {
		return Decimal{}, fmt.Errorf("not a valid decimal number: %s", input)
	}

	value := shift(x, digits) + y
	if negative {
		value = -value
	}
	return Decimal{Value: value, Digits: digits}, nil
}

func parseInt(input string) (result, digits int, negative, ok bool) {
	runes := []rune(input)

	if len(runes) > 0 && runes[0] == '-' {
		negative = true
		runes = runes[1:]
	}

	for i, digit := range runes {
		if i != 0 {
			result *= 10
		}
		if digit < '0' || digit > '9' {
			return
		}
		value := int(digit - '0')
		result += value
		digits++
	}

	ok = true
	return
}

func shift(x, digits int) int {
	for i := 0; i < digits; i++ {
		x *= 10
	}
	return x
}
