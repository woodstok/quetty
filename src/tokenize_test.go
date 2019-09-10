package quetty

import (
	"testing"
)

func TestTokenize(t *testing.T) {
	basicString := "This is a sentence"
	basicStringExpected := []string{"This", "is", "a", "sentence"}
	t.Run("Nil Tokenizer", func(t *testing.T) {
		got, err := Tokenize(basicString, &NilTokenizer{})
		assertNoError(t, err)
		assertSliceEqual(t, []string{basicString}, got)
	})

	t.Run("tokenize basic set of strings", func(t *testing.T) {
		got, err := Tokenize(basicString, &BasicTokenizer{})
		assertNoError(t, err)
		assertSliceEqual(t, basicStringExpected, got)
	})

	regexTests := []struct {
		name    string
		pattern string
		input   string
		err     bool
		want    []string
	}{
		{name: "Word", pattern: WORDREGEX, input: basicString, err: false, want: basicStringExpected},
		{name: "Number", pattern: NUMREGEX, input: "2 numbers in 486 sentences", err: false, want: []string{"2", "486"}},
		{name: "Hash", pattern: HASHREGEX, input: "4c 2abd4c 4beefc1 abcdefghijk abcdef9gh", err: false, want: []string{"2abd4c", "4beefc1"}},
		{name: "BadRegex", pattern: `[a-f0-9A-F`, input: "Does not matter", err: true, want: nil},
	}

	for _, tt := range regexTests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Tokenize(tt.input, &RegexTokenizer{pattern: tt.pattern})
			if tt.err {
				assertSomeError(t, err)
			} else {
				assertNoError(t, err)
			}
			assertSliceEqual(t, tt.want, got)
		})

	}

}
