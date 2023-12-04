package main

import (
	"fmt"
	"os"
	"path"
	"regexp"
	"strconv"
	"sync"
)

type Ref struct {
	S    Symbol
	X, Y int
}

type Symbol string

func (s Symbol) IsProperSymbol() bool {
	return !regexp.MustCompile(`^([.|\d])$`).MatchString(string(s))
}

func (s Symbol) IsDigit() bool {
	return regexp.MustCompile(`^(\d)$`).MatchString(string(s))
}

type Grid map[int]map[int]Symbol

func (g *Grid) Walk(cb func(Symbol, int, int) bool) []*Ref {
	var mu sync.Mutex
	var refs []*Ref

	var wg sync.WaitGroup
	wg.Add(len(*g) * len((*g)[1]))
	go func() {
		for y, row := range *g {
			for x, val := range row {
				go func(val Symbol, x, y int) {
					ref := &Ref{val, x, y}
					defer wg.Done()
					if cb(val, x, y) {
						mu.Lock()
						defer mu.Unlock()
						refs = append(refs, ref)
					}
				}(val, x, y)
			}
		}
	}()
	wg.Wait()

	return refs
}

func (g *Grid) Get(x, y int) *Symbol {
	if col, cok := (*g)[y]; cok {
		if val, rok := col[x]; rok {
			return &val
		}
	}
	return nil
}

func ReadGrid(c string) Grid {
	g := Grid{}
	y, x := 1, 1
	for i := 0; i < len(c); i++ {
		ch := fmt.Sprintf("%c", c[i])
		if ch == "\n" {
			y += 1
			x = 1
			continue
		}

		if v, ok := g[y]; !ok {
			g[y] = map[int]Symbol{
				x: Symbol(ch),
			}
		} else {
			v[x] = Symbol(ch)
		}

		x += 1
	}
	return g
}

func First(c []byte) (int, error) {
	g := ReadGrid(string(c))

	refs := g.Walk(func(val Symbol, x, y int) bool {
		if !regexp.MustCompile(`^\d$`).MatchString(string(val)) {
			return false
		}

		tl := g.Get(x-1, y-1)
		t := g.Get(x, y-1)
		tr := g.Get(x+1, y-1)
		r := g.Get(x+1, y)
		rb := g.Get(x+1, y+1)
		b := g.Get(x, y+1)
		bl := g.Get(x-1, y+1)
		l := g.Get(x-1, y)

		checkRel := func(s *Symbol) bool {
			if s == nil {
				return false
			}

			if !s.IsProperSymbol() {
				return false
			}

			return true
		}

		if checkRel(tl) ||
			checkRel(t) ||
			checkRel(tr) ||
			checkRel(r) ||
			checkRel(rb) ||
			checkRel(b) ||
			checkRel(bl) ||
			checkRel(l) {
			return true
		}

		return false
	})

	res := 0
	lookup := map[string]bool{}

	for _, r := range refs {
		sign := fmt.Sprintf("%vx%v", r.X, r.Y)
		if v, ok := lookup[sign]; ok && v == true {
			continue
		}
		lookup[sign] = true

		before := []string{}
		for i := 1; ; i += 1 {
			if v := g.Get(r.X-i, r.Y); v == nil {
				break
			} else if v.IsDigit() {
				lookup[fmt.Sprintf("%vx%v", r.X-i, r.Y)] = true
				before = append(before, string(*v))
			} else {
				break
			}
		}

		after := []string{}
		for i := 1; ; i += 1 {
			if v := g.Get(r.X+i, r.Y); v == nil {
				break
			} else if v.IsDigit() {
				lookup[fmt.Sprintf("%vx%v", r.X+i, r.Y)] = true
				after = append(after, string(*v))
			} else {
				break
			}
		}

		s := ""
		if len(before) > 0 {
			for i := len(before) - 1; i >= 0; i-- {
				s += before[i]
			}
		}
		s += string(r.S)
		for _, a := range after {
			s += a
		}

		n, err := strconv.Atoi(s)
		if err != nil {
			return 0, err
		}

		res += n
	}

	return res, nil
}

func Second(c []byte) (int, error) {
	g := ReadGrid(string(c))

	refs := g.Walk(func(val Symbol, x, y int) bool {
		if !val.IsDigit() {
			return false
		}

		tl := g.Get(x-1, y-1)
		t := g.Get(x, y-1)
		tr := g.Get(x+1, y-1)
		r := g.Get(x+1, y)
		rb := g.Get(x+1, y+1)
		b := g.Get(x, y+1)
		bl := g.Get(x-1, y+1)
		l := g.Get(x-1, y)

		checkRel := func(s *Symbol) bool {
			if s == nil {
				return false
			}

			if *s == "*" {
				return true
			}

			return false
		}

		if checkRel(tl) ||
			checkRel(t) ||
			checkRel(tr) ||
			checkRel(r) ||
			checkRel(rb) ||
			checkRel(b) ||
			checkRel(bl) ||
			checkRel(l) {
			return true
		}

		return false
	})

	cogs := g.Walk(func(val Symbol, x, y int) bool {
		return val == "*"
	})

	res := 0
	lookup := map[string]bool{}

	type combination struct {
		number int
		coords []struct {
			x, y int
		}
	}

	combos := []combination{}

	for _, r := range refs {
		sign := fmt.Sprintf("%vx%v", r.X, r.Y)
		if v, ok := lookup[sign]; ok && v == true {
			continue
		}
		lookup[sign] = true

		combo := combination{
			coords: []struct {
				x, y int
			}{
				{x: r.X, y: r.Y},
			},
		}

		before := []string{}
		for i := 1; ; i += 1 {
			if v := g.Get(r.X-i, r.Y); v == nil {
				break
			} else if v.IsDigit() {
				lookup[fmt.Sprintf("%vx%v", r.X-i, r.Y)] = true
				before = append(before, string(*v))
				combo.coords = append(combo.coords, struct {
					x, y int
				}{r.X - i, r.Y})
			} else {
				break
			}
		}

		after := []string{}
		for i := 1; ; i += 1 {
			if v := g.Get(r.X+i, r.Y); v == nil {
				break
			} else if v.IsDigit() {
				lookup[fmt.Sprintf("%vx%v", r.X+i, r.Y)] = true
				after = append(after, string(*v))
				combo.coords = append(combo.coords, struct {
					x, y int
				}{r.X + i, r.Y})
			} else {
				break
			}
		}

		s := ""
		if len(before) > 0 {
			for i := len(before) - 1; i >= 0; i-- {
				s += before[i]
			}
		}
		s += string(r.S)
		for _, a := range after {
			s += a
		}

		n, err := strconv.Atoi(s)
		if err != nil {
			return 0, err
		}

		combo.number = n
		combos = append(combos, combo)
	}

	type bundle struct {
		cog     *Ref
		numbers []int
	}

	var bundles []bundle
	for _, c := range cogs {
		b := bundle{
			cog:     c,
			numbers: []int{},
		}

		for _, combo := range combos {
			matched := false
			for _, coords := range combo.coords {
				if coords.x >= c.X-1 &&
					coords.x <= c.X+1 &&
					coords.y <= c.Y+1 &&
					coords.y >= c.Y-1 {
					if matched {
						break
					}
					matched = true
					b.numbers = append(b.numbers, combo.number)
				}
			}
		}

		bundles = append(bundles, b)
	}

	for _, b := range bundles {
		if len(b.numbers) != 2 {
			continue
		}

		v := 1
		for _, n := range b.numbers {
			v *= n
		}
		res += v
	}

	return res, nil
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

	res, err := First(c)
	if err != nil {
		panic(err)
	}
	fmt.Println("#1: ", res)

	con, err := os.ReadFile(path.Join(cwd, "2-input.txt"))
	if err != nil {
		panic(err)
	}

	sum, err := Second(con)
	if err != nil {
		panic(err)
	}
	fmt.Println("#2: ", sum)
}
