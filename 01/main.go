package main

import (
	"fmt"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
)

func First(lines []string) (int, error) {
	sum := 0
	first := regexp.MustCompile(`^[^\d]*(\d)`)
	last := regexp.MustCompile(`(\d)[^\d]*$`)
	for _, l := range lines {
		mf := first.FindStringSubmatch(l)
		ml := last.FindStringSubmatch(l)
		r := mf[1] + ml[1]

		lineCount, err := strconv.Atoi(r)
		if err != nil {
			return 0, err
		}

		sum += lineCount
	}

	return sum, nil
}

func Second(lines []string) (int, error) {
	mapping := []struct {
		Word, Replacement string
	}{
		{"one", "1"},
		{"two", "2"},
		{"three", "3"},
		{"four", "4"},
		{"five", "5"},
		{"six", "6"},
		{"seven", "7"},
		{"eight", "8"},
		{"nine", "9"},
	}

	names := []string{}
	for _, m := range mapping {
		names = append(names, m.Word)
	}

	reverse := func(s string) string {
		letters := strings.Split(s, "")
		n := ""
		for l := len(letters) - 1; l >= 0; l -= 1 {
			n += letters[l]
		}
		return n
	}

	reverseNames := []string{}
	for _, m := range mapping {
		reverseNames = append(reverseNames, reverse(m.Word))
	}

	lookup := map[string]string{}
	for _, m := range mapping {
		lookup[m.Word] = m.Replacement
	}

	reverseLookup := map[string]string{}
	for _, m := range mapping {
		reverseLookup[reverse(m.Word)] = m.Replacement
	}

	nameRe := regexp.MustCompile("(" + strings.Join(names, "|") + ")")
	reverseNameRe := regexp.MustCompile("(" + strings.Join(reverseNames, "|") + ")")

	var newLines []string
	for _, l := range lines {
		newLine := l
		forwardCalled := false
		newLine = nameRe.ReplaceAllStringFunc(newLine, func(s string) string {
			if forwardCalled {
				return s
			}

			forwardCalled = true
			return lookup[s]
		})

		reverseCalled := false
		reverseLine := reverse(newLine)
		reverseLine = reverseNameRe.ReplaceAllStringFunc(reverseLine, func(s string) string {
			if reverseCalled {
				return s
			}

			reverseCalled = true
			return reverseLookup[s]
		})

		newLine = reverse(reverseLine)
		newLines = append(newLines, newLine)
	}

	return First(newLines)
}

func main() {
	cwd, _ := os.Getwd()
	getLines := func(path string) []string {
		contents, err := os.ReadFile(path)
		if err != nil {
			panic(err)
		}
		return strings.Split(string(contents), "\n")
	}

	_, err := First(getLines(path.Join(cwd, "01-input.txt")))
	if err != nil {
		panic(err)
	}

	//fmt.Printf("#1 result: %v\n", first)

	second, err := Second(getLines(path.Join(cwd, "01-input.txt")))
	if err != nil {
		panic(err)
	}

	fmt.Printf("#2 result: %v\n", second)
}
