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

func TestIpTokenizer(t *testing.T) {
	IpTestString := `
	10.0.0.1
	192.168.1.1
	2001:0db8:85a3:0000:0000:8a2e:0370:7334
	2001:db8:85a3::8a2e:370:7334
	`
	IpTestExpected := []string{
		"10.0.0.1",
		"192.168.1.1",
		"2001:0db8:85a3:0000:0000:8a2e:0370:7334",
		"2001:db8:85a3::8a2e:370:7334",
	}
	t.Run("tokenize ip addresses", func(t *testing.T) {
		got, err := Tokenize(IpTestString, &IpTokenizer{})
		assertNoError(t, err)
		assertSliceEqual(t, IpTestExpected, got)
	})

}

func TestPathTokenizer(t *testing.T) {
	pathTestString := `
	file1
	file2.txt
	file3.tar.gz
	/file4
	/file5.jpg
	/file6.tgz
	/root/path1/path2
	/root/path3/path4/
	/root/path5/path6/filename.ext
	/root/path7/path8/filename.ext.gz
	root/path9/path10
	root/path11/path14/
	root/path13/path16/filename.ext
	root/path15/path18/filename.ext.gz
	root/pat-h.9/path10
	root/pat-h.11/path14/
	root/pat-h.13/path16/filename.ext
	root/pat-h.15/path18/filename.ext.gz
	https://testurl
	http://testurl/1.txt
	10.0.0.1
	192.168.1.1
	`
	pathTestExpected := []string{

		// "file1",
		"file2.txt",
		"file3.tar.gz",
		"/file4",
		"/file5.jpg",
		"/file6.tgz",
		"/root/path1/path2",
		"/root/path3/path4/",
		"/root/path5/path6/filename.ext",
		"/root/path7/path8/filename.ext.gz",
		"root/path9/path10",
		"root/path11/path14/",
		"root/path13/path16/filename.ext",
		"root/path15/path18/filename.ext.gz",
		"root/pat-h.9/path10",
		"root/pat-h.11/path14/",
		"root/pat-h.13/path16/filename.ext",
		"root/pat-h.15/path18/filename.ext.gz",
	}
	t.Run("tokenize string of paths", func(t *testing.T) {
		got, err := Tokenize(pathTestString, &PathTokenizer{})
		assertNoError(t, err)
		assertSliceEqual(t, pathTestExpected, got)
	})

}
