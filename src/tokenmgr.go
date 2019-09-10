package quetty

import (
	"bytes"
	"fmt"
	"io"
)

type TokenMgr struct {
	options    *Options
	tokenizers []Tokenizer
}

func NewTokenMgr(options *Options) (*TokenMgr, error) {
	tokMgr := TokenMgr{options: options}

	if options.matchWord {
		tokMgr.tokenizers = append(tokMgr.tokenizers,
			&RegexTokenizer{pattern: WORDREGEX})
	}
	if options.matchNum {
		tokMgr.tokenizers = append(tokMgr.tokenizers,
			&RegexTokenizer{pattern: NUMREGEX})
	}
	if options.matchHash {
		tokMgr.tokenizers = append(tokMgr.tokenizers,
			&RegexTokenizer{pattern: HASHREGEX})
	}

	if len(tokMgr.tokenizers) == 0 {
		return nil, fmt.Errorf("No tokenizers specified")
	}

	return &tokMgr, nil
}

func (tMgr *TokenMgr) Process(reader io.Reader) (Tokens, error) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(reader)
	inputString := buf.String()
	tokens := NewTokens(nil)
	for _, t := range tMgr.tokenizers {
		tokenSlice, err := Tokenize(inputString, t)
		if err != nil {
			return nil, err
		}
		tokens.Extend(tokenSlice)
	}
	return tokens, nil
}
