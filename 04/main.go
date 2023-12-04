package main

import (
	"fmt"
	"os"
	"path"
	"regexp"
	"strings"
)

func First(lines []string) (int, error) {
	sum := 0
	for _, l := range lines {
		card := strings.Split(l, ":")
		numbers := strings.Split(card[1], "|")
		winning := regexp.MustCompile(`(\d+)`).FindAllStringSubmatch(numbers[0], -1)

		winningNumbers := []string{}
		for _, w := range winning {
			winningNumbers = append(winningNumbers, w[1])
		}

		matches := 0
		winningRe := regexp.MustCompile(fmt.Sprintf(`^(%s)$`, strings.Join(winningNumbers, "|")))
		for _, s := range strings.Split(numbers[1], " ") {
			if winningRe.MatchString(s) {
				matches += 1
			}
		}

		n := 0
		if matches > 0 {
			n = 1
		}

		for t := 1; t < matches; t++ {
			n = n * 2
		}

		sum += n
	}

	return sum, nil
}

func Second(lines []string) (int, error) {
	cards := map[int]int{}
	for idx, _ := range lines {
		cards[idx+1] = 1
	}

	for lid, l := range lines {
		idx := lid + 1
		card := strings.Split(l, ":")
		numbers := strings.Split(card[1], "|")
		winning := regexp.MustCompile(`(\d+)`).FindAllStringSubmatch(numbers[0], -1)

		winningNumbers := []string{}
		for _, w := range winning {
			winningNumbers = append(winningNumbers, w[1])
		}

		matches := 0
		winningRe := regexp.MustCompile(fmt.Sprintf(`^(%s)$`, strings.Join(winningNumbers, "|")))
		for _, s := range strings.Split(numbers[1], " ") {
			if winningRe.MatchString(s) {
				matches += 1
			}
		}

		for w := 1; w <= matches; w++ {
			v := cards[idx+w]
			v += 1 * (cards[idx])
			cards[idx+w] = v
		}
	}

	sum := 0
	for _, a := range cards {
		sum += a
	}

	return sum, nil
}

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	c, err := os.ReadFile(path.Join(cwd, "1-input.txt"))
	if err != nil {
		panic(err)
	}

	lines := strings.Split(string(c), "\n")
	sum, err := First(lines)
	if err != nil {
		panic(err)
	}

	fmt.Println("#1: ", sum)

	total, err := Second(lines)
	if err != nil {
		panic(err)
	}

	fmt.Println("#2: ", total)
}
