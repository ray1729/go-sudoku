package main

import (
	"fmt"
	"strings"
)

type IntSet map[int]struct{}

func NewIntSet(xs []int) IntSet {
	s := make(IntSet)
	for _, x := range xs {
		s[x] = struct{}{}
	}
	return s
}

func (s IntSet) Add(i int) IntSet {
	_, ok := s[i]
	if !ok {
		s[i] = struct{}{}
	}
	return s
}

func (s IntSet) Remove(i int) IntSet {
	if _, ok := s[i]; ok {
		delete(s, i)
	}
	return s
}

func (s IntSet) Contains(i int) bool {
	_, ok := s[i]
	return ok
}

func (s IntSet) Elements() []int {
	elements := make([]int, 0, len(s))
	for v := range s {
		elements = append(elements, v)
	}
	return elements
}

type Grid [81]int

func (g Grid) At(row, col int) int {
	return g[row*9+col]
}

func (g Grid) String() string {
	var sb strings.Builder
	for r := 0; r < 9; r++ {
		fmt.Fprintf(&sb, "%d %d %d | %d %d %d | %d %d %d\n",
			g.At(r, 0), g.At(r, 1), g.At(r, 2),
			g.At(r, 3), g.At(r, 4), g.At(r, 5),
			g.At(r, 6), g.At(r, 7), g.At(r, 8))
		if r == 2 || r == 5 {
			sb.WriteString("------+-------+------\n")
		}
	}
	return sb.String()
}

func (g Grid) Neighbours(i int) IntSet {
	// Grid is stored in row-major order, so i = row * 9 + col
	col := i % 9
	row := (i - col) / 9
	// Top-left corner of the minor square
	tlRow := (row / 3) * 3
	tlCol := (col / 3) * 3

	neighbours := make(IntSet)

	// Neighbouring row
	for c := 0; c < 9; c++ {
		if v := g.At(row, c); v != 0 {
			neighbours.Add(v)
		}
	}

	// Neighbouring col
	for r := 0; r < 9; r++ {
		if v := g.At(r, col); v != 0 {
			neighbours.Add(v)
		}
	}

	// Minor square
	for r := tlRow; r < tlRow+3; r++ {
		for c := tlCol; c < tlCol+3; c++ {
			if v := g.At(r, c); v != 0 {
				neighbours.Add(v)
			}
		}
	}

	return neighbours
}

func (g Grid) FirstUnsolved() int {
	for i, v := range g {
		if v == 0 {
			return i
		}
	}
	return -1
}

func (g Grid) WithElementAt(i int, v int) Grid {
	g[i] = v
	return g
}

func (g Grid) Solve() (Grid, error) {
	frontier := []Grid{g}
	for len(frontier) > 0 {
		n := len(frontier) - 1
		g := frontier[n]
		frontier = frontier[:n]
		i := g.FirstUnsolved()
		if i < 0 {
			// We found a solution
			return g, nil
		}
		neighbours := g.Neighbours(i)
		for v := 1; v < 10; v++ {
			if neighbours.Contains(v) {
				continue
			}
			frontier = append(frontier, g.WithElementAt(i, v))
		}
	}
	// We have explored all frontiers without finding a solution
	return Grid{}, fmt.Errorf("No solutions found")
}

func main() {
	G := Grid{
		5, 3, 0, 0, 7, 0, 0, 0, 0,
		6, 0, 0, 1, 9, 5, 0, 0, 0,
		0, 9, 8, 0, 0, 0, 0, 6, 0,
		8, 0, 0, 0, 6, 0, 0, 0, 3,
		4, 0, 0, 8, 0, 3, 0, 0, 1,
		7, 0, 0, 0, 2, 0, 0, 0, 6,
		0, 6, 0, 0, 0, 0, 2, 8, 0,
		0, 0, 0, 4, 1, 9, 0, 0, 5,
		0, 0, 0, 0, 8, 0, 0, 7, 9,
	}

	solution, err := G.Solve()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(solution)
	}
}
