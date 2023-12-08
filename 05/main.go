package main

import (
	"fmt"
	"math"
	"os"
	"path"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
)

type Range struct {
	Start, Length int
}

func (r Range) End() int {
	return r.Start + (r.Length - 1)
}

func (r Range) Overlap(cr Range) []int {
	if (cr.Start >= r.Start && cr.Start <= r.End()) ||
		(cr.End() >= r.Start && cr.End() <= r.End()) ||
		(r.Start >= cr.Start && r.Start <= cr.End()) ||
		(r.End() >= cr.Start && r.End() <= cr.End()) {
		start := math.Max(float64(cr.Start), float64(r.Start))
		end := math.Min(float64(cr.End()), float64(r.End()))

		var n []int
		for s := start; s <= end; s++ {
			n = append(n, int(s))
		}

		return n
	}

	return []int{}
}

type RangeAssignment struct {
	SourceRange, DestinationRange Range
}

func (ra *RangeAssignment) MapNumber(in int) (int, bool) {
	if in >= ra.SourceRange.Start && in <= ra.SourceRange.End() {
		offset := in - ra.SourceRange.Start
		out := ra.DestinationRange.Start + offset
		return out, true
	}

	return in, false
}

func InterpretMapLine(def string) (*RangeAssignment, error) {
	n := regexp.MustCompile(`^(\d+)\s+(\d+)\s+(\d+)$`).FindStringSubmatch(def)
	destStart, err := strconv.Atoi(n[1])
	if err != nil {
		return nil, err
	}

	srcStart, err := strconv.Atoi(n[2])
	if err != nil {
		return nil, err
	}

	distance, err := strconv.Atoi(n[3])
	if err != nil {
		return nil, err
	}

	return &RangeAssignment{
		SourceRange: Range{
			Start:  srcStart,
			Length: distance,
		},
		DestinationRange: Range{
			Start:  destStart,
			Length: distance,
		},
	}, nil
}

func ProcessStage(mapLines []string) func(int) (int, error) {
	return func(in int) (int, error) {
		for _, l := range mapLines {
			ra, err := InterpretMapLine(l)
			if err != nil {
				return 0, err
			}

			if out, m := ra.MapNumber(in); m {
				return out, nil
			}
		}

		return in, nil
	}
}

func ExtractSeeds(l string) []int {
	numbers := strings.Split(l, ":")[1]
	n := regexp.MustCompile(`(\d+)`).FindAllStringSubmatch(numbers, -1)
	var seeds []int
	for _, m := range n {
		s, err := strconv.Atoi(m[1])
		if err != nil {
			panic(err)
		}
		seeds = append(seeds, s)
	}
	return seeds
}

func First(seeds []int, stages [][]string) (int, error) {
	funnel := []func(int) (int, error){}
	for _, s := range stages {
		funnel = append(funnel, ProcessStage(s))
	}

	var wg sync.WaitGroup
	wg.Add(len(seeds))

	errs := make(chan error, len(seeds))
	res := make(chan int, len(seeds))
	for _, s := range seeds {
		go func(in int) {
			var out int
			var err error
			for _, f := range funnel {
				out, err = f(in)
				if err != nil {
					errs <- err
					break
				}
				in = out
			}
			res <- out
			wg.Done()
		}(s)
	}
	wg.Wait()
	close(res)

	var locs []int
	for l := range res {
		locs = append(locs, l)
	}

	sort.Ints(locs)
	return locs[0], nil
}

// todo(troublete): doesn't work for real input; too slow –> need to implement reverse lookup for best path
func Second(seeds []int, stages [][]string) (int, error) {
	funnel := []func(int) (int, error){}

	var seedRanges []Range
	for cur := 0; len(seeds)-2 >= cur; cur += 2 {
		seedRanges = append(seedRanges, Range{
			Start:  seeds[cur],
			Length: seeds[cur+1],
		})
	}

	var firstStages []Range
	lines := stages[0]
	for _, l := range lines {
		a, err := InterpretMapLine(l)
		if err != nil {
			return 0, err
		}

		firstStages = append(firstStages, a.SourceRange)
	}

	var processable []int
	for _, sr := range seedRanges {
		for _, stage := range firstStages {
			processable = append(processable, sr.Overlap(stage)...)
		}
	}

	for _, s := range stages {
		funnel = append(funnel, ProcessStage(s))
	}

	var wg sync.WaitGroup
	wg.Add(len(processable))

	errs := make(chan error, len(processable))
	res := make(chan int, len(processable))
	for _, s := range processable {
		go func(in int) {
			var out int
			var err error
			for _, f := range funnel {
				out, err = f(in)
				if err != nil {
					errs <- err
					break
				}
				in = out
			}
			res <- out
			wg.Done()
		}(s)
	}
	wg.Wait()
	close(res)

	var locs []int
	for l := range res {
		locs = append(locs, l)
	}

	sort.Ints(locs)
	if len(locs) > 0 {
		return locs[0], nil
	}
	return 0, nil
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

	var seeds []int
	var stages [][]string
	var stage []string

	lines := strings.Split(string(c), "\n")
	for _, l := range lines {
		if strings.Trim(l, " ") == "" {
			continue
		}

		if regexp.MustCompile(`^seeds`).MatchString(l) {
			seeds = append(seeds, ExtractSeeds(l)...)
			continue
		}

		if regexp.MustCompile(`map:`).MatchString(l) {
			if len(stage) > 0 {
				stages = append(stages, stage)
				stage = []string{}
			}
			continue
		}

		stage = append(stage, l)
	}

	stages = append(stages, stage)

	result, err := First(seeds, stages)
	if err != nil {
		panic(err)
	}

	fmt.Println("#1:", result)

	second, err := Second(seeds, stages)
	if err != nil {
		panic(err)
	}

	fmt.Println("#2:", second)
}
