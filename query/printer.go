package query

import (
	"fmt"
	"strings"
)

type Printer struct {
	builder     strings.Builder
	indentation int
}

func (p *Printer) Print(format string, a ...any) {
	for i := 0; i < p.indentation; i++ {
		fmt.Fprint(&p.builder, "    ")
	}
	fmt.Fprintf(&p.builder, format, a...)
	fmt.Fprintln(&p.builder)
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
