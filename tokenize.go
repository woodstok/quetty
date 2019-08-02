package main

import "strings"

type Tokenizer func(input string) []string

func Tokenize(input string, tokFunc Tokenizer) []string {
	return tokFunc(input)
}

func BasicTokenizer(input string) []string {
	return strings.Split(input, " ")
}
