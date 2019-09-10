package quetty

import (
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/andreyvit/diff"
)

func assertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("didn't expect an error but got one, %v", err)
	}
}

func assertSomeError(t *testing.T, got error) {
	t.Helper()
	if got == nil {
		t.Fatal("did not get an error but wanted one", got)
	}
}

func assertError(t *testing.T, got error, want error) {
	t.Helper()
	if got == nil {
		t.Fatal("didn't get an error but wanted one")
	}

	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}

func assertSliceEqual(t *testing.T, want, got []string) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("want %v, got %v", want, got)
	}
}

func assertTokensEqual(t *testing.T, want, got Tokens) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("want %v, got %v", want, got)
	}

}

func AssertSortedEqual(t *testing.T, testName, want, got string) {
	t.Helper()
	expectedLines := strings.Split(want, "\n")
	gotLines := strings.Split(got, "\n")
	sort.Strings(expectedLines)
	sort.Strings(gotLines)
	if !reflect.DeepEqual(expectedLines, gotLines) {
		sortedExpected := strings.Join(expectedLines, "\n")
		sortedGot := strings.Join(gotLines, "\n")
		t.Errorf("Ouput not as expected for test '%s'\n%v", testName, diff.LineDiff(sortedExpected, sortedGot))

	}

}
