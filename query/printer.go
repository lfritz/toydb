package query

import (
	"fmt"
	"strings"
)

type Printer struct {
	builder     strings.Builder
	indentation int
	inLine      bool
}

func (p *Printer) Println(format string, a ...any) {
	p.Print(format, a...)
	fmt.Fprintln(&p.builder)
	p.inLine = false
}

func (p *Printer) Print(format string, a ...any) {
	if !p.inLine {
		for i := 0; i < p.indentation; i++ {
			fmt.Fprint(&p.builder, "    ")
		}
	}
	fmt.Fprintf(&p.builder, format, a...)
	p.inLine = true
}

func (p *Printer) Indent() {
	p.indentation++
}

func (p *Printer) Unindent() {
	if p.indentation > 0 {
		p.indentation--
	}
}

func (p *Printer) String() string {
	return p.builder.String()
}

type Printable interface {
	Print(printer *Printer)
}

func Print(p Printable) string {
	printer := new(Printer)
	p.Print(printer)
	return printer.String()
}
