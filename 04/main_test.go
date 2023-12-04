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

	c, err := os.ReadFile(path.Join(cwd, "1-test-input.txt"))
	if err != nil {
		t.Error(err)
	}

	lines := strings.Split(string(c), "\n")
	sum, err := First(lines)
	if err != nil {
		t.Error(err)
	}

	want := 13
	if sum != want {
		t.Errorf("want %v, got %v", want, sum)
	}
}

func Test_Second(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Error(err)
	}

	c, err := os.ReadFile(path.Join(cwd, "2-test-input.txt"))
	if err != nil {
		t.Error(err)
	}

	lines := strings.Split(string(c), "\n")
	sum, err := Second(lines)
	if err != nil {
		t.Error(err)
	}

	want := 30
	if sum != want {
		t.Errorf("want %v, got %v", want, sum)
	}
}
