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

func getLines(path string) []string {
	f, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	lines := strings.Split(string(f), "\n")
	return lines
}

type Config struct {
	Reds, Greens, Blues int
}

func (cfg Config) Power() float64 {
	r := math.Max(1.0, float64(cfg.Reds))
	g := math.Max(1.0, float64(cfg.Greens))
	b := math.Max(1.0, float64(cfg.Blues))

	return r * g * b
}

type Game struct {
	ID     int
	Rounds []Config
}

func (g Game) Valid(cfg Config) bool {
	for _, r := range g.Rounds {
		if r.Reds > cfg.Reds || r.Blues > cfg.Blues || r.Greens > cfg.Greens {
			return false
		}
	}
	return true
}

func (g Game) SmallestConfigPossible() Config {
	cfg := Config{}
	for _, r := range g.Rounds {
		if cfg.Reds < r.Reds {
			cfg.Reds = r.Reds
		}
		if cfg.Blues < r.Blues {
			cfg.Blues = r.Blues
		}
		if cfg.Greens < r.Greens {
			cfg.Greens = r.Greens
		}
	}
	return cfg
}

type GameList []Game

func (gl GameList) ValidSum(cfg Config) int {
	sum := 0

	for _, g := range gl {
		if g.Valid(cfg) {
			sum += g.ID
		}
	}

	return sum
}

func First(lines []string) (*GameList, error) {
	var gl GameList
	for _, l := range lines {
		parts := strings.Split(l, ":")
		idx := regexp.MustCompile(`^Game\s(\d+)`).FindStringSubmatch(parts[0])

		gameID, err := strconv.Atoi(idx[1])
		if err != nil {
			return nil, err
		}

		game := Game{ID: gameID, Rounds: []Config{}}
		rounds := strings.Split(parts[1], ";")
		for _, r := range rounds {
			cfg := Config{}
			colors := strings.Split(r, ",")
			for _, c := range colors {
				m := regexp.MustCompile(`(\d+)\s+((red|green|blue))`).FindStringSubmatch(c)
				amount, err := strconv.Atoi(m[1])
				if err != nil {
					return nil, err
				}
				switch m[2] {
				case "red":
					cfg.Reds = amount
				case "green":
					cfg.Greens = amount
				case "blue":
					cfg.Blues = amount
				}
			}
			game.Rounds = append(game.Rounds, cfg)
		}
		gl = append(gl, game)
	}

	return &gl, nil
}

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	result, err := First(getLines(path.Join(cwd, "01-input.txt")))
	if err != nil {
		panic(err)
	}

	fmt.Println("#1:", result.ValidSum(Config{
		Reds:   12,
		Greens: 13,
		Blues:  14,
	}))

	power := 0.0
	for _, g := range *result {
		power += g.SmallestConfigPossible().Power()
	}

	fmt.Println("#2:", power)
}
