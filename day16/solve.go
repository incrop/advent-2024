package day16

import (
	"container/heap"
	"strconv"
)

type Solve struct{}

func (d Solve) Part1(input []string, outputCh chan<- []string) string {
	f := parse(input)
	outputCh <- f.output()
	exitState, _ := f.cheapestPathToExit()
	outputCh <- f.output()
	return strconv.FormatInt(exitState.score, 10)
}

func (d Solve) Part2(input []string, outputCh chan<- []string) string {
	f := parse(input)
	outputCh <- f.output()
	exitState, visited := f.cheapestPathToExit()
	outputCh <- f.output()
	bestTilesCount := f.bestTilesCount(exitState, visited)
	outputCh <- f.output()
	return strconv.Itoa(bestTilesCount)
}

type coord [2]int

func (c1 coord) plus(c2 coord) coord {
	return coord{c1[0] + c2[0], c1[1] + c2[1]}
}

type field struct {
	cells [][]rune
	start coord
	exit  coord
}

func (f field) cell(coord coord) *rune {
	return &f.cells[coord[0]][coord[1]]
}

type dir int

var dirRunes = []rune("→↓←↑")
var dirOffsets [4]coord = [4]coord{{0, 1}, {1, 0}, {0, -1}, {-1, 0}}

type location struct {
	coord coord
	dir   dir
}

type state struct {
	location
	score int64
}

type minScoreHeap []state

func (h minScoreHeap) Len() int           { return len(h) }
func (h minScoreHeap) Less(i, j int) bool { return h[i].score < h[j].score }
func (h minScoreHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *minScoreHeap) Push(x any) {
	*h = append(*h, x.(state))
}
func (h *minScoreHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

const (
	stepCost = 1
	turnCost = 1000
)

func (f field) cheapestPathToExit() (exitState state, visited map[location]int64) {
	minScores := minScoreHeap{}
	heap.Push(&minScores, state{location{f.start, 0}, 0})
	visited = map[location]int64{}
	for {
		s := heap.Pop(&minScores).(state)
		if s.coord == f.exit {
			return s, visited
		}
		if _, ok := visited[s.location]; ok {
			continue
		}
		visited[s.location] = s.score
		r := f.cell(s.coord)
		if *r == '.' {
			*r = dirRunes[s.dir]
		}
		step := location{s.coord.plus(dirOffsets[s.dir]), s.dir}
		if _, ok := visited[step]; *f.cell(step.coord) != '#' && !ok {
			heap.Push(&minScores, state{step, s.score + stepCost})
		}
		for _, turn := range [2]location{{s.coord, (s.dir + 1) % 4}, {s.coord, (s.dir + 3) % 4}} {
			nextStepCoord := turn.coord.plus(dirOffsets[turn.dir])
			if _, ok := visited[step]; *f.cell(nextStepCoord) != '#' && !ok {
				heap.Push(&minScores, state{turn, s.score + turnCost})
			}
		}
	}
}

func (f field) bestTilesCount(exitState state, visited map[location]int64) (count int) {
	best := map[location]bool{}
	var backtrack = func(curr state) {}
	backtrack = func(curr state) {
		best[curr.location] = true
		if r := f.cell(curr.coord); *r != '_' {
			*r = '_'
			count++
		}
		for _, prev := range []state{
			{
				location{
					curr.coord.plus(dirOffsets[(curr.dir+2)%4]),
					curr.dir,
				},
				curr.score - stepCost,
			},
			{
				location{
					curr.coord,
					(curr.dir + 3) % 4,
				},
				curr.score - turnCost,
			},
			{
				location{
					curr.coord,
					(curr.dir + 1) % 4,
				},
				curr.score - turnCost,
			},
		} {
			if score, ok := visited[prev.location]; ok && score == prev.score && !best[prev.location] {
				backtrack(prev)
			}
		}
	}
	backtrack(exitState)
	return
}

func parse(input []string) (f field) {
	f.cells = make([][]rune, len(input))
	for i, line := range input {
		f.cells[i] = []rune(line)
		for j, r := range line {
			switch r {
			case 'S':
				f.start = [2]int{i, j}
			case 'E':
				f.exit = [2]int{i, j}
			}
		}
	}
	return
}

func (f field) output() (lines []string) {
	for _, row := range f.cells {
		lines = append(lines, string(row))
	}
	return
}

func (d Solve) CorrectAnswers() [2]string {
	return [2]string{"72400", "435"}
}
