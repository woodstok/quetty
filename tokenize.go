package main

import "strings"

type Tokenizer interface {
	Tokenize(input string) []string
}

func Tokenize(input string, tok Tokenizer) []string {
	return tok.Tokenize(input)
}

type BasicTokenizer struct {
}

func (t *BasicTokenizer) Tokenize(input string) []string {
	return strings.Split(input, " ")
}

type NilTokenizer struct {
}

func (t *NilTokenizer) Tokenize(input string) []string {
	return []string{input}
}
