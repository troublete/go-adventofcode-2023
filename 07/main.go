package main

import (
	"fmt"
	"os"
	"path"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

var (
	Values = map[string]int{
		"A": 14,
		"K": 13,
		"Q": 12,
		"J": 11,
		"T": 10,
	}
)

type Card struct {
	Label string
	Value int
}

func NewCard(l string) (Card, error) {
	var v int
	var ok bool
	var err error
	if v, ok = Values[l]; !ok {
		v, err = strconv.Atoi(l)
	}

	if err != nil {
		return Card{}, err
	}
	return Card{
		Label: l,
		Value: v,
	}, nil
}

type Hand struct {
	Cards  map[int]Card
	Makeup map[string][]*Card
	Bid    int
}

func (h Hand) Value() int {
	switch len(h.Makeup) {
	case 1: // five of a kind
		return 6
	case 2:
		for _, cs := range h.Makeup {
			if len(cs) == 4 {
				return 5 // four of a kind
			}

			if len(cs) == 2 || len(cs) == 3 {
				return 4 // full house
			}
		}
	case 3:
		for _, cs := range h.Makeup {
			if len(cs) == 3 {
				return 3 // triplet
			}
		}

		fallthrough
	case 4:
		nPairs := 0
		for _, cs := range h.Makeup {
			if len(cs) == 2 {
				nPairs += 1
			}
		}

		if nPairs > 1 {
			return 2 // two pair
		}

		if nPairs == 1 {
			return 1 // single pair
		}
	default:
		return 0
	}
	return 0
}

type Game []Hand

func (g Game) Order() {
	sort.Slice(g, func(a, b int) bool {
		aHand := g[a]
		bHand := g[b]

		aVal := aHand.Value()
		bVal := bHand.Value()

		switch {
		case aVal > bVal:
			return true
		case bVal > aVal:
			return false
		case aVal == bVal:
			for idx, aC := range aHand.Cards {
				switch {
				case aC.Value > bHand.Cards[idx].Value:
					return true
				case aC.Value < bHand.Cards[idx].Value:
					return false
				}
			}
		}

		return true
	})
}

func (g Game) Sum() int {
	g.Order()

	sum := 0
	for idx, h := range g {
		for _, c := range h.Cards {
			fmt.Print(c.Label)
		}

		fmt.Print(" ")
		fmt.Print(h.Bid)
		fmt.Print("\n")

		sum += h.Bid * (len(g) - idx)
	}

	return sum
}

func First(lines []string) (int, error) {
	var hands Game
	for _, l := range lines {
		m := regexp.MustCompile(`^(.)(.)(.)(.)(.)\s(\d+)`).FindStringSubmatch(l)
		bid, err := strconv.Atoi(m[6])
		if err != nil {
			return 0, err
		}
		h := Hand{
			Cards:  map[int]Card{},
			Makeup: map[string][]*Card{},
			Bid:    bid,
		}
		for idx, c := range m[1:6] {
			card, err := NewCard(c)
			if err != nil {
				return 0, err
			}
			h.Cards[idx+1] = card

			if v, ok := h.Makeup[card.Label]; !ok {
				h.Makeup[card.Label] = []*Card{&card}
			} else {
				l := v
				l = append(l, &card)
				h.Makeup[card.Label] = l
			}
		}
		hands = append(hands, h)
	}

	return hands.Sum(), nil
}

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	contents, err := os.ReadFile(path.Join(cwd, "1-input.txt"))
	if err != nil {
		panic(err)
	}

	lines := strings.Split(string(contents), "\n")
	_, err = First(lines)
	if err != nil {
		panic(err)
	}

	// fmt.Println("#1:", first)
}
