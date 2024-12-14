package day10

import "strings"

type Calculate struct{}

func (d Calculate) Part1(input []string, outputCh chan<- []string) int64 {
	f := parse(input)
	outputCh <- f.output()
	return int64(f.calculateScoreSum())
}

func (f field) calculateScoreSum() (sum int) {
	cache := make([][]map[[2]int]bool, len(f))
	for i, row := range f {
		cache[i] = make([]map[[2]int]bool, len(row))
	}
	var collectReachableCoords func(i, j int) map[[2]int]bool
	collectReachableCoords = func(i, j int) map[[2]int]bool {
		cell := f[i][j]
		if cell.value == '9' {
			return map[[2]int]bool{{i, j}: true}
		}
		if cached := cache[i][j]; cached != nil {
			return cached
		}
		coordsUnion := make(map[[2]int]bool)
		for dir, off := range dirOffsets {
			if !cell.conn[dir] {
				continue
			}
			for coord := range collectReachableCoords(i+off[0], j+off[1]) {
				coordsUnion[coord] = true
			}
		}
		cache[i][j] = coordsUnion
		return coordsUnion
	}
	for i, row := range f {
		for j, cell := range row {
			if cell.value == '0' {
				sum += len(collectReachableCoords(i, j))
			}
		}
	}
	return
}

func (d Calculate) Part2(input []string, outputCh chan<- []string) int64 {
	f := parse(input)
	outputCh <- f.output()
	return int64(f.calculateRatingSum())
}

func (f field) calculateRatingSum() (sum int64) {
	cache := make([][]int64, len(f))
	for i, row := range f {
		cache[i] = make([]int64, len(row))
	}
	var calculateRecursively func(i, j int) int64
	calculateRecursively = func(i, j int) int64 {
		cell := f[i][j]
		if cell.value == '9' {
			return 1
		}
		if cached := cache[i][j]; cached > 0 {
			return cached
		}
		dirSum := int64(0)
		for dir, off := range dirOffsets {
			if !cell.conn[dir] {
				continue
			}
			dirSum += calculateRecursively(i+off[0], j+off[1])
		}
		cache[i][j] = dirSum
		return dirSum
	}
	for i, row := range f {
		for j, cell := range row {
			if cell.value == '0' {
				sum += calculateRecursively(i, j)
			}
		}
	}
	return
}

var dirOffsets [4][2]int = [4][2]int{{-1, 0}, {0, 1}, {1, 0}, {0, -1}}

type cell struct {
	value rune
	conn  [4]bool
}

type field [][]cell

func parse(input []string) (f field) {
	for i, line := range input {
		row := make([]cell, len(line))
		for j, r := range line {
			cell := &row[j]
			cell.value = r
			for dir, off := range dirOffsets {
				ii, jj := i+off[0], j+off[1]
				if ii < 0 || ii >= len(input) || jj < 0 || jj >= len(line) {
					continue
				}
				rr := input[ii][jj]
				if r+1 != rune(rr) {
					continue
				}
				cell.conn[dir] = true
			}
		}
		f = append(f, row)
	}
	return
}

func (f field) output() (lines []string) {
	for i, row := range f {
		var sb strings.Builder
		for j, cell := range row {
			if cell.value == '0' || cell.value == '9' {
				sb.WriteRune(cell.value)
				continue
			}
			biConn := cell.conn
			for dir, off := range dirOffsets {
				ii, jj := i+off[0], j+off[1]
				if ii < 0 || ii >= len(f) || jj < 0 || jj >= len(row) {
					continue
				}
				if f[ii][jj].conn[(dir+2)%4] {
					biConn[dir] = true
				}
			}
			switch biConn {
			case [4]bool{false, false, true, true}:
				sb.WriteRune('╗')
			case [4]bool{false, true, false, true}:
				sb.WriteRune('═')
			case [4]bool{false, true, true, false}:
				sb.WriteRune('╔')
			case [4]bool{false, true, true, true}:
				sb.WriteRune('╦')
			case [4]bool{true, false, false, true}:
				sb.WriteRune('╝')
			case [4]bool{true, false, true, false}:
				sb.WriteRune('║')
			case [4]bool{true, false, true, true}:
				sb.WriteRune('╣')
			case [4]bool{true, true, false, false}:
				sb.WriteRune('╚')
			case [4]bool{true, true, false, true}:
				sb.WriteRune('╩')
			case [4]bool{true, true, true, false}:
				sb.WriteRune('╠')
			case [4]bool{true, true, true, true}:
				sb.WriteRune('╬')
			default:
				sb.WriteRune(cell.value)
			}
		}
		lines = append(lines, sb.String())
	}
	return
}

func (d Calculate) CorrectAnswers() [2]int64 {
	return [2]int64{698, 1436}
}
