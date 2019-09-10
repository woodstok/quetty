package quetty

import (
	"fmt"
	"sort"
	"strings"
)

type void struct{}

var member void

type Tokens map[string]void

func NewTokens(list []string) Tokens {
	tokens := make(Tokens)
	tokens.Extend(list)
	return tokens
}

func (t Tokens) Extend(list []string) {
	for _, val := range list {
		t.Add(val)
	}
}

func (t Tokens) Add(val string) {
	t[val] = member
}

func (t Tokens) Del(val string) {
	delete(t, val)
}

func (t Tokens) Contains(val string) bool {
	_, exists := t[val]
	return exists
}

func (t Tokens) ToSlice() []string {
	ret := make([]string, 0, len(t))
	for val, _ := range t {
		ret = append(ret, val)
	}
	return ret
}
func (t Tokens) ToSortedSlice() []string {
	ret := make([]string, 0, len(t))
	for val, _ := range t {
		ret = append(ret, val)
	}
	sort.Strings(ret)
	return ret
}

func (t Tokens) String() string {
	var sb strings.Builder
	for tok, _ := range t {
		sb.WriteString(tok)
	}
	return fmt.Sprintf("Tokens{%s}", strings.Join(t.ToSlice(), ", "))
}
