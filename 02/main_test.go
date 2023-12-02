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

	res, err := First(getLines(path.Join(cwd, "01-test-input.txt")))

	if err != nil {
		t.Error(err)
	}

	want := 8
	if res.ValidSum(Config{
		Reds:   12,
		Greens: 13,
		Blues:  14,
	}) != want {
		t.Errorf("want %v, got %v", want, res)
	}
}

func Test_Second(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Error(err)
	}

	res, err := First(getLines(path.Join(cwd, "01-test-input.txt")))

	if err != nil {
		t.Error(err)
	}

	power := 0.0
	for _, g := range *res {
		power += g.SmallestConfigPossible().Power()
	}

	want := 2286.0
	if power != want {
		t.Errorf("want %v, got %v", want, res)
	}
}
