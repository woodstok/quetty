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
	basicString := "This is a sentence"
	t.Run("Nil Tokenizer", func(t *testing.T) {
		got := Tokenize(basicString, &NilTokenizer{})
		assertSliceEqual(t, []string{basicString}, got)
	})

	t.Run("tokenize basic set of strings", func(t *testing.T) {
		want := []string{"This", "is", "a", "sentence"}
		got := Tokenize(basicString, &BasicTokenizer{})
		assertSliceEqual(t, want, got)
	})
}
