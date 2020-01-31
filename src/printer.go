package quetty

import (
	"fmt"
	"io"
)

type Printer struct {
	options *Options
}

func NewPrinter(opts *Options) (*Printer, error) {
	return &Printer{options: opts}, nil
}

func (p *Printer) Print(writer io.WriteCloser, tokens Tokens) {
	for token, _ := range tokens {
		fmt.Fprintln(writer, token)
	}
	writer.Close()
}
