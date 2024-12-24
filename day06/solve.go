package day06

import (
	"strconv"
	"strings"
)

type Solve struct{}

func (d Solve) Part1(input []string, outputCh chan<- []string) string {
	f := parse(input)
	outputCh <- f.output()
	visitedCount := 0
	for f.isInside(f.pos) {
		nextPos := f.nextPos()
		if f.isInside(nextPos) && f.cell(nextPos).value != '.' {
			f.dir = (f.dir + 1) % 4
		} else {
			currCell := f.cell(f.pos)
			if currCell.visited == [4]bool{} {
				visitedCount++
			}
			currCell.visited[f.dir] = true
			f.pos = nextPos
		}
	}
	outputCh <- f.output()
	return strconv.Itoa(visitedCount)
}

func (d Solve) Part2(input []string, outputCh chan<- []string) string {
	f := parse(input)
	found := make(map[[2]int]bool)
	f.findLoopPositionsRecusive(found, true, outputCh)
	return strconv.Itoa(len(found))
}

func (f field) findLoopPositionsRecusive(found map[[2]int]bool, blockLeft bool, outputCh chan<- []string) bool {
	currCell := f.cell(f.pos)
	currVisited := &currCell.visited[f.dir]
	if *currVisited {
		outputCh <- f.output()
		return true
	}
	*currVisited = true
	defer func() {
		*currVisited = false
	}()

	nextPos := f.nextPos()
	if !f.isInside(nextPos) {
		return false
	}
	if f.cell(nextPos).value != '.' {
		f.dir = (f.dir + 1) % 4
		return f.findLoopPositionsRecusive(found, blockLeft, outputCh)
	}

	if nextCell := f.cell(nextPos); blockLeft && nextCell.visited == [4]bool{} {
		nextCell.value = 'O'
		f.dir = (f.dir + 1) % 4
		if f.findLoopPositionsRecusive(found, false, outputCh) {
			found[nextPos] = true
		}
		f.dir = (f.dir + 3) % 4
		nextCell.value = '.'
	}

	f.pos = nextPos
	return f.findLoopPositionsRecusive(found, blockLeft, outputCh)
}

type dir int

const (
	up dir = iota
	right
	down
	left
)

type cell struct {
	value   rune
	visited [4]bool
}

type field struct {
	cells [][]cell
	pos   [2]int
	dir   dir
}

func parse(input []string) (f field) {
	for i, line := range input {
		row := make([]cell, len(line))
		for j, c := range line {
			if c == '^' {
				f.pos[0] = i
				f.pos[1] = j
				row[j].value = '.'
			} else {
				row[j].value = c
			}
		}
		f.cells = append(f.cells, row)
	}
	return
}

func (f field) isInside(pos [2]int) bool {
	if pos[0] < 0 || pos[0] >= len(f.cells) {
		return false
	}
	row := f.cells[pos[0]]
	if pos[1] < 0 || pos[1] >= len(row) {
		return false
	}
	return true
}

func (f field) cell(pos [2]int) *cell {
	return &f.cells[pos[0]][pos[1]]
}

func (f field) nextPos() [2]int {
	np := f.pos
	switch f.dir {
	case up:
		np[0]--
	case right:
		np[1]++
	case down:
		np[0]++
	case left:
		np[1]--
	}
	return np
}

func (f field) output() (o []string) {
	if f.isInside(f.pos) {
		posCell := f.cell(f.pos)
		saved := posCell.value
		posCell.value = rune("^>v<"[f.dir])
		defer func() {
			posCell.value = saved
		}()
	}
	for _, row := range f.cells {
		var sb strings.Builder
		for _, cell := range row {
			if cell.value == '.' && cell.visited != [4]bool{} {
				sb.WriteRune('+')
			} else {
				sb.WriteRune(cell.value)
			}
		}
		o = append(o, sb.String())
	}
	return
}

func (d Solve) CorrectAnswers() [2]string {
	return [2]string{"5067", "1793"}
}
