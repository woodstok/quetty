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
	t.Run("tokenize basic set of strings", func(t *testing.T) {
		want := []string{"This", "is", "a", "sentence"}
		got := Tokenize("This is a sentence")
		assertSliceEqual(t, want, got)
	})
}
