package main

import (
	quetty "github.com/woodstok/quetty/src"
)

func main() {
	quetty.Tokenize("testing", &quetty.NilTokenizer{})
}
