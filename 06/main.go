package main

import (
	"fmt"
	"math"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
)

func CalculateWinningWays(time, winDistance int) int {
	var ways int

	for i := time; i > time/2; i-- {
		if i*(time-i) >= winDistance {
			ways = int(math.Max(float64(i), float64(time-i))-math.Min(float64(i), float64(time-i))) + 1
			break
		}
	}

	return ways
}

func First(lines []string) (int, error) {
	var times []int
	var distances []int

	for _, l := range lines {
		if regexp.MustCompile(`^Time:`).MatchString(l) {
			n := regexp.MustCompile(`(\d+)`).FindAllStringSubmatch(strings.Split(l, ":")[1], -1)
			for _, m := range n {
				time, err := strconv.Atoi(m[1])
				if err != nil {
					return 0, err
				}

				times = append(times, time)
			}
			continue
		}

		if regexp.MustCompile(`^Distance:`).MatchString(l) {
			n := regexp.MustCompile(`(\d+)`).FindAllStringSubmatch(strings.Split(l, ":")[1], -1)
			for _, m := range n {
				distance, err := strconv.Atoi(m[1])
				if err != nil {
					return 0, err
				}

				distances = append(distances, distance)
			}
		}
	}

	sum := 1
	for idx, t := range times {
		sum *= CalculateWinningWays(t, distances[idx]+1)
	}

	return sum, nil
}

func Second(lines []string) (int, error) {
	var time int
	var distance int
	var err error

	for _, l := range lines {
		if regexp.MustCompile(`^Time:`).MatchString(l) {
			n := strings.Split(l, ":")[1]
			time, err = strconv.Atoi(strings.ReplaceAll(n, " ", ""))
			if err != nil {
				return 0, err
			}
		}

		if regexp.MustCompile(`^Distance:`).MatchString(l) {
			n := strings.Split(l, ":")[1]
			distance, err = strconv.Atoi(strings.ReplaceAll(n, " ", ""))
			if err != nil {
				return 0, err
			}
		}
	}

	return CalculateWinningWays(time, distance), nil
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
	first, err := First(lines)
	if err != nil {
		panic(err)
	}
	fmt.Println("#1:", first)

	second, err := Second(lines)
	if err != nil {
		panic(err)
	}
	fmt.Println("#2:", second)
}
