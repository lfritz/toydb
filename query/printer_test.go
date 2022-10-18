package query

import "testing"

func TestPrinter(t *testing.T) {
	printer := new(Printer)
	printer.Println("foo")
	printer.Indent()
	printer.Print("bar: ")
	printer.Println("baz")
	printer.Unindent()
	printer.Println("qux")
	got := printer.String()
	want := `foo
    bar: baz
qux
`
	if got != want {
		t.Errorf("printer produced\n%s\nwant\n%s", got, want)
	}
}
