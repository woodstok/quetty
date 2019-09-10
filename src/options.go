package quetty

import (
	"flag"
	"fmt"
)

type Options struct {
	matchWord bool
	matchNum  bool
	matchHash bool
	minLen    uint
}

func initFlags(opt *Options) {
	flag.BoolVar(&opt.matchWord, "word", false, "Tokenize basic words(w+)")
	flag.BoolVar(&opt.matchNum, "num", false, "Tokenize basic numbers")
	flag.BoolVar(&opt.matchHash, "hash", false, "Tokenize hash values")
	flag.UintVar(&opt.minLen, "m", 4, "minimum length of tokens")
}

func ParseOptions() *Options {
	opts := &Options{}
	initFlags(opts)
	flag.Parse()
	fmt.Printf("opts = %+v\n", opts)
	return opts
}
