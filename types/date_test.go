package types

import "testing"

func TestCheckDateValid(t *testing.T) {
	cases := []struct {
		year, month, day int
		want             Date
	}{
		{1999, 1, 31, Date{1999, 1, 31}},
		{2000, 2, 29, Date{2000, 2, 29}},
		{1, 1, 1, Date{1, 1, 1}},
		{9999, 12, 31, Date{9999, 12, 31}},
	}
	for _, c := range cases {
		got, ok := CheckDate(c.year, c.month, c.day)
		if !ok {
			t.Errorf("CheckDate(%v, %v, %v) returned ok = false", c.year, c.month, c.day)
		}
		if got != c.want {
			t.Errorf("CheckDate(%v, %v, %v) returned %v, want %v", c.year, c.month, c.day, got, c.want)
		}
	}
}

func TestCheckDateInvalid(t *testing.T) {
	cases := []struct {
		year, month, day int
	}{
		{0, 12, 31},
		{10000, 1, 1},
		{1999, 12, 0},
		{1999, 12, 32},
		{1999, 2, 29},
	}
	for _, c := range cases {
		_, ok := CheckDate(c.year, c.month, c.day)
		if ok {
			t.Errorf("CheckDate(%v, %v, %v) returned ok = true", c.year, c.month, c.day)
		}
	}
}

func TestParseDate(t *testing.T) {
	input := "1999-12-31"
	got, err := ParseDate(input)
	if err != nil {
		t.Errorf("ParseDate(%q) returned error: %v", input, err)
	}
	want := NewDate(1999, 12, 31)
	if got != want {
		t.Errorf("ParseDate(%q) returned %v, want %v", input, got, want)
	}

	invalid := []string{
		"1999-12-1",
		"1999-12-31 ",
		" 1999-12-31",
		"1999-12-32",
	}
	for _, input := range invalid {
		_, err := ParseDate(input)
		if err == nil {
			t.Errorf("ParseDate(%q) did not return error", input)
		}
	}
}

func TestDaysInMonth(t *testing.T) {
	cases := []struct {
		year, month int
		want        int
	}{
		{1999, 1, 31},
		{1999, 2, 28},
		{1999, 4, 30},
		{2000, 2, 29},
	}
	for _, c := range cases {
		got := daysInMonth(c.year, c.month)
		if got != c.want {
			t.Errorf("daysInMonth(%v, %v) == %v, want %v", c.year, c.month, got, c.want)
		}
	}
}

func TestLeapYear(t *testing.T) {
	cases := []struct {
		year int
		want bool
	}{
		{1999, false},
		{1996, true},
		{1900, false},
		{2000, true},
	}
	for _, c := range cases {
		got := leapYear(c.year)
		if got != c.want {
			t.Errorf("leapYear(%v) == %v, want %v", c.year, got, c.want)
		}
	}
}

func TestDateCompare(t *testing.T) {
	cases := []struct {
		a, b Date
		want Compared
	}{
		{Date{2000, 1, 1}, Date{1999, 12, 31}, ComparedGt},
		{Date{1999, 12, 31}, Date{2000, 1, 1}, ComparedLt},
		{Date{1999, 12, 31}, Date{1999, 12, 31}, ComparedEq},
	}
	for _, c := range cases {
		got := c.a.Compare(c.b)
		if got != c.want {
			t.Errorf("%v.Compare(%v) == %v, want %v", c.a, c.b, got, c.want)
		}
	}
}
