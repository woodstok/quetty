package main

import (
	"reflect"
	"testing"
)

func assertSliceEqual(t *testing.T, want, got []string) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("want %v, got %v", want, got)
	}

}

func TestTokenize(t *testing.T) {
	t.Run("Nil Tokenizer", func(t *testing.T) {
		basicString := "This is a sentence"
		got := Tokenize(basicString, func(input string) []string {
			return []string{input}

		})
		assertSliceEqual(t, []string{basicString}, got)
	})

	t.Run("tokenize basic set of strings", func(t *testing.T) {
		want := []string{"This", "is", "a", "sentence"}
		got := Tokenize("This is a sentence", BasicTokenizer)
		assertSliceEqual(t, want, got)
	})
}
