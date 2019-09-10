package quetty

import (
	"strings"
	"testing"
)

func TestTokenMgr(t *testing.T) {
	basicString := "This is a sentence 43 567123 abcd dcba dead123beef non123hash"
	wantNum := "43 567123 123"
	wantHash := "abcd dcba dead123beef"
	wantNumAndHash := "43 567123 123 abcd dcba dead123beef"
	wantMin4Basic := "This sentence 567123 abcd dcba dead123beef non123hash"
	wantMin4Num := "567123"
	wantMin4Hash := "abcd dcba dead123beef"
	wantMin4NumAndHash := "567123 abcd dcba dead123beef"
	t.Run("token manager will return error on no tokenizers", func(t *testing.T) {
		_, err := NewTokenMgr(&Options{})
		assertSomeError(t, err)
	})

	tokenMgrTests := []struct {
		name      string
		options   *Options
		inputStr  string
		expectErr bool
		wantStr   string
	}{
		{name: "tokenize basic words", options: &Options{matchWord: true},
			inputStr: basicString, expectErr: false, wantStr: basicString},
		{name: "tokenize numbers", options: &Options{matchNum: true},
			inputStr: basicString, expectErr: false, wantStr: wantNum},
		{name: "tokenize hashes", options: &Options{matchHash: true},
			inputStr: basicString, expectErr: false, wantStr: wantHash},
		{name: "tokenize hashes and numbers", options: &Options{matchHash: true, matchNum: true},
			inputStr: basicString, expectErr: false, wantStr: wantNumAndHash},
		{name: "tokenize basic words with 4 minimum", options: &Options{minLen: 4, matchWord: true},
			inputStr: basicString, expectErr: false, wantStr: wantMin4Basic},
		{name: "tokenize numbers with 4 minimum", options: &Options{minLen: 4, matchNum: true},
			inputStr: basicString, expectErr: false, wantStr: wantMin4Num},
		{name: "tokenize hashes with 4 minimum", options: &Options{minLen: 4, matchHash: true},
			inputStr: basicString, expectErr: false, wantStr: wantMin4Hash},
		{name: "tokenize hashes and numbers with 4 minimum", options: &Options{minLen: 4, matchHash: true, matchNum: true},
			inputStr: basicString, expectErr: false, wantStr: wantMin4NumAndHash},
	}

	for _, tt := range tokenMgrTests {
		t.Run(tt.name, func(t *testing.T) {
			tokMgr, err := NewTokenMgr(tt.options)
			assertNoError(t, err)

			reader := strings.NewReader(tt.inputStr)
			tokens, err := tokMgr.Process(reader)
			assertNoError(t, err)

			want := NewTokens(strings.Split(tt.wantStr, " "))
			assertTokensEqual(t, want, tokens)
		})

	}
}
