package quetty

import (
	"regexp"
	"strings"
)

const (
	WORDREGEX = `\w+`
	NUMREGEX  = `\d+`

	// technically all numbers match a basic hashregex
	// but let us only match hashes that has
	// atleast one alphabet in it
	HASHREGEX = `([a-f0-9A-F]*[a-fA-F][a-f0-9A-F]*){4,}\b`
)

type Tokenizer interface {
	Tokenize(input string) ([]string, error)
}

func Tokenize(input string, tok Tokenizer) ([]string, error) {
	return tok.Tokenize(input)
}

// Basic Tokenizer
// splits based on space
type BasicTokenizer struct {
}

func (t *BasicTokenizer) Tokenize(input string) ([]string, error) {
	return strings.Split(input, " "), nil
}

// Nil Tokenizer
// Returns the input as is
type NilTokenizer struct {
}

func (t *NilTokenizer) Tokenize(input string) ([]string, error) {
	return []string{input}, nil
}

// Regex Tokenizer
// splits based on space
type RegexTokenizer struct {
	pattern string
	re_     *regexp.Regexp
}

func (t *RegexTokenizer) Tokenize(input string) ([]string, error) {
	var err error
	if t.re_ == nil {
		t.re_, err = regexp.Compile(t.pattern)
		if err != nil {
			return nil, err
		}
	}
	return t.re_.FindAllString(input, -1), nil
}
