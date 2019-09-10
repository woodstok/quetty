package quetty

import (
	"strings"
	"testing"
)

func TestTokenMgr(t *testing.T) {
	basicString := "This is a sentence 43 567 abcd dcba dead123beef non123hash"
	t.Run("token manager will return error on no tokenizers", func(t *testing.T) {
		_, err := NewTokenMgr(&Options{})
		assertSomeError(t, err)
	})
	t.Run("basic word tokenset", func(t *testing.T) {
		tokMgr, err := NewTokenMgr(&Options{matchWord: true})
		assertNoError(t, err)

		reader := strings.NewReader(basicString)
		tokens, err := tokMgr.Process(reader)
		assertNoError(t, err)

		want := NewTokens(strings.Split(basicString, " "))
		assertTokensEqual(t, want, tokens)

	})
}
