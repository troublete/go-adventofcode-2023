package main

import (
	"os"
	"path"
	"strings"
	"testing"
)

func Test_First(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Error(err)
	}

	contents, err := os.ReadFile(path.Join(cwd, "01-test-input.txt"))
	if err != nil {
		t.Error(err)
	}

	lines := strings.Split(string(contents), "\n")
	result, err := First(lines)
	if err != nil {
		t.Error(err)
	}

	want := 142
	if result != want {
		t.Errorf("want %v, got %v", want, result)
	}
}

func Test_Second(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Error(err)
	}

	contents, err := os.ReadFile(path.Join(cwd, "02-test-input.txt"))
	if err != nil {
		t.Error(err)
	}

	lines := strings.Split(string(contents), "\n")
	result, err := Second(lines)
	if err != nil {
		t.Error(err)
	}

	want := 281
	if result != want {
		t.Errorf("want %v, got %v", want, result)
	}
}
