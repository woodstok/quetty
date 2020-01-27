package quetty

import (
	"fmt"
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

func assertStringEqual(t *testing.T, want, got string) {
	t.Helper()
	if got != want {
		t.Errorf("want %s\ngot %s\n", want, got)
	}
}

func assertSliceEqual(t *testing.T, want, got []string) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("want %v\ngot %v\ndiff %s\n", want, got, sliceDiff(want, got))
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

func sliceDiff(want, got []string) string {
	sort.Strings(want)
	sort.Strings(got)
	sortedExpected := strings.Join(want, "\n")
	sortedGot := strings.Join(got, "\n")
	return diff.LineDiff(sortedExpected, sortedGot)
}

func sliceDifference(want, got []string) string {
	var sb strings.Builder
	sort.Strings(want)
	sort.Strings(got)
	wLen := len(want)
	gLen := len(got)
	i := 0
	j := 0
	sb.WriteString("[")

	for i < wLen && j < gLen {
		if want[i] == got[j] {
			i++
			j++
		} else if want[i] < got[j] {
			sb.WriteString(fmt.Sprintf("-%s ", want[i]))
			i++

		} else {
			sb.WriteString(fmt.Sprintf("+%s ", got[j]))
			j++
		}
	}

	for i < wLen {
		sb.WriteString(fmt.Sprintf("-%s ", want[i]))
		i++
	}
	for j < gLen {
		sb.WriteString(fmt.Sprintf("+%s ", got[i]))
		j++

	}
	sb.WriteString("]")
	return sb.String()
}
