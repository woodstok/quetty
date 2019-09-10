package quetty

import (
	"log"
	"os"
)

func Run(opts *Options) {

	tokenMgr, err := NewTokenMgr(opts)
	if err != nil {
		log.Fatal(err)
	}

	tokens, err := tokenMgr.Process(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	printer, err := NewPrinter(opts)
	if err != nil {
		log.Fatal(err)
	}
	printer.Print(os.Stdout, tokens)
}
