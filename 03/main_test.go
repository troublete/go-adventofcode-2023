package main

import (
	"os"
	"path"
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

	res, err := First(c)
	if err != nil {
		t.Error(err)
	}

	want := 4361
	if res != want {
		t.Errorf("want %v, got %v", want, res)
	}
}

func Test_Second(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Error(err)
	}

	c, err := os.ReadFile(path.Join(cwd, "1-test-input.txt"))
	if err != nil {
		t.Error(err)
	}

	res, err := Second(c)
	if err != nil {
		t.Error(err)
	}

	want := 467835
	if res != want {
		t.Errorf("want %v, got %v", want, res)
	}
}
