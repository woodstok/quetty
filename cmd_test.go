package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"testing"

	quetty "github.com/woodstok/quetty/src"
)

var update = flag.Bool("update", false, "update golden files")

var binaryName = "bin/quetty"

func fixturePath(t *testing.T, fixture string) string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatalf("problems recovering caller information")
	}

	return filepath.Join(filepath.Dir(filename), "golden", fixture)
}

func writeFixture(t *testing.T, fixture string, content []byte) {
	err := ioutil.WriteFile(fixturePath(t, fixture), content, 0644)
	if err != nil {
		t.Fatal(err)
	}
}

func loadFixture(t *testing.T, fixture string) string {
	content, err := ioutil.ReadFile(fixturePath(t, fixture))
	if err != nil {
		t.Fatal(err)
	}

	return string(content)
}

func TestCliArgs(t *testing.T) {
	tests := []struct {
		name          string
		args          []string
		inputFixture  string
		outputFixture string
	}{
		{"words", []string{"-word"}, "words", "words"},
		{"numbers", []string{"-num"}, "words", "num"},
		{"hash", []string{"-hash"}, "hash", "hash"},
		{"numhash", []string{"-hash", "-num"}, "hash", "numhash"},
		{"path", []string{"-path"}, "path", "path"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir, err := os.Getwd()
			if err != nil {
				t.Fatal(err)
			}

			cmd := exec.Command(path.Join(dir, binaryName), tt.args...)

			inputFile := tt.inputFixture + "-input"
			cmd.Stdin, err = os.Open(fixturePath(t, inputFile))
			if err != nil {
				t.Fatal(err)
			}

			output, err := cmd.CombinedOutput()
			if err != nil {
				t.Fatal(err)
			}

			outputFile := tt.outputFixture + "-output"
			if *update {
				writeFixture(t, outputFile, output)
			}

			actual := string(output)

			expected := loadFixture(t, outputFile)

			quetty.AssertSortedEqual(t, tt.name, expected, actual)
		})
	}
}

func TestMain(m *testing.M) {
	make := exec.Command("make")
	err := make.Run()
	if err != nil {
		fmt.Printf("could not make binary for %s: %v\n", binaryName, err)
		os.Exit(1)
	}

	os.Exit(m.Run())
}
