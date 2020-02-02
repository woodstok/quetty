package quetty

import (
	"flag"
)

type Options struct {
	init         bool
	stdin        bool
	matchWord    bool
	matchNum     bool
	matchHash    bool
	matchPath    bool
	matchIp      bool
	matchTime    bool
	matchNospace bool
	matchIdent   bool
	minLen       uint
}

func initFlags(opt *Options) {
	flag.BoolVar(&opt.init, "init", false, "Start fzf with quetty keybindings")
	flag.BoolVar(&opt.stdin, "stdin", false, "Basic tokenizer read from stdin. Don't start fzf interface")
	flag.BoolVar(&opt.matchWord, "word", false, "Tokenize basic words(w+)")
	flag.BoolVar(&opt.matchNum, "num", false, "Tokenize basic numbers")
	flag.BoolVar(&opt.matchHash, "hash", false, "Tokenize hash values")
	flag.BoolVar(&opt.matchPath, "path", false, "Tokenize filepaths")
	flag.BoolVar(&opt.matchIp, "ip", false, "Tokenize ip addresses")
	flag.BoolVar(&opt.matchTime, "time", false, "Tokenize time addresses")
	flag.BoolVar(&opt.matchNospace, "nospace", false, "Tokenize all nonspace tokens")
	flag.BoolVar(&opt.matchIdent, "ident", false, "Tokenize identifiers")
	flag.UintVar(&opt.minLen, "m", 4, "minimum length of tokens")
}

func ParseOptions() *Options {
	opts := &Options{}
	initFlags(opts)
	flag.Parse()
	// fmt.Printf("opts = %+v\n", opts)
	return opts
}
