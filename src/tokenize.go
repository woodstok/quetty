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

// Path Tokenizer
// should contain atleast one back slash
// or be file with a dot ( for files like 'filename.tar.gz'
// will ignore plain pathnames like 'filename' for now
const (
	// pathRegex = `\b(http(s)?://)?([A-Za-z_0-9@~:.-/]*/[A-Za-z_0-9@~:.-]*)\b`
	pathRegex            = `((http(s)?://)?([A-Za-z_0-9@~:./-]*/[A-Za-z_0-9@~:./-]*))`
	fileNameWithDotRegex = `(\b[A-Za-z0-9@:-]+\.[A-Za-z0-9@:.-]+\b)`
)

type PathTokenizer struct {
	re_ *regexp.Regexp
}

func (t *PathTokenizer) isValid(token string) bool {
	if isMatch, _ := regexp.MatchString(`[a-zA-Z]`, token); !isMatch {
		return false
	}
	if strings.HasPrefix(token, "http") {
		return false
	}
	return true

}

func (t *PathTokenizer) validStrings(tokens []string) []string {
	ret := make([]string, 0)
	for _, tok := range tokens {
		if t.isValid(tok) {
			ret = append(ret, tok)
		}
	}
	return ret

}
func (t *PathTokenizer) Tokenize(input string) ([]string, error) {
	var err error
	if t.re_ == nil {
		combinedRegex := "(" + pathRegex + "|" + fileNameWithDotRegex + ")"
		t.re_, err = regexp.Compile(combinedRegex)
		if err != nil {
			return nil, err
		}
	}
	return t.validStrings(t.re_.FindAllString(input, -1)), nil
}
