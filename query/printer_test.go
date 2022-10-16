package query

import "testing"

func TestPrinter(t *testing.T) {
	printer := new(Printer)
	printer.Print("foo")
	printer.Indent()
	printer.Print("bar")
	printer.Unindent()
	printer.Print("baz")
	got := printer.String()
	want := `foo
    bar
baz
`
	if got != want {
		t.Errorf("printer produced\n%s\nwant\n%s", got, want)
	}
}
