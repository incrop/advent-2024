package day04

import "strconv"

type Solve struct{}

func (d Solve) Part1(input []string, outputCh chan<- []string) string {
	xmasCount := int64(0)
	f := fieldFrom(input)
	for i, line := range input {
		for j := range line {
			dirs := xmasDirections(input, i, j)
			xmasCount += int64(len(dirs))
			for _, dir := range dirs {
				f.copy(4, i, j, dir[0], dir[1])
				outputCh <- f.output()
			}
		}
	}
	return strconv.FormatInt(xmasCount, 10)
}

func xmasDirections(input []string, i, j int) (dirs [][2]int) {
	for di := -1; di <= 1; di++ {
		for dj := -1; dj <= 1; dj++ {
			if di == 0 && dj == 0 {
				continue
			}
			if readsXmsInThisDirection(input, i, j, di, dj) {
				dirs = append(dirs, [2]int{di, dj})
			}
		}
	}
	return
}

func readsXmsInThisDirection(input []string, i, j, di, dj int) bool {
	for _, xmasLetter := range "XMAS" {
		if i < 0 || i >= len(input) {
			return false
		}
		line := input[i]
		if j < 0 || j >= len(line) {
			return false
		}
		if line[j] != byte(xmasLetter) {
			return false
		}
		i += di
		j += dj
	}
	return true
}

func (d Solve) Part2(input []string, outputCh chan<- []string) string {
	xmasCount := int64(0)
	f := fieldFrom(input)
	for i, line := range input {
		for j := range line {
			if hasMasInXShapeAt(input, i, j) {
				xmasCount++
				f.copy(3, i-1, j-1, 1, 1)
				f.copy(3, i-1, j+1, 1, -1)
				outputCh <- f.output()
			}
		}
	}
	outputCh <- f.output()
	return strconv.FormatInt(xmasCount, 10)
}

func hasMasInXShapeAt(input []string, i, j int) bool {
	if input[i][j] != 'A' {
		return false
	}
	for di, dj := -1, -1; di <= 1; di += 2 {
		msCount := [2]int{}
		for _, coord := range [2][2]int{{i + di, j + dj}, {i - di, j - dj}} {
			if coord[0] < 0 || coord[0] >= len(input) {
				return false
			}
			line := input[coord[0]]
			if coord[1] < 0 || coord[1] >= len(line) {
				return false
			}
			switch line[coord[1]] {
			case 'M':
				msCount[0]++
			case 'S':
				msCount[1]++
			default:
				return false
			}
		}
		if msCount != [2]int{1, 1} {
			return false
		}
	}
	return true
}

type field struct {
	input []string
	cells [][]rune
}

func fieldFrom(input []string) (f field) {
	f.input = input
	for _, line := range input {
		row := make([]rune, len(line))
		for j := range row {
			row[j] = ' '
		}
		f.cells = append(f.cells, row)
	}
	return
}

func (f field) copy(n, i, j, di, dj int) {
	for range n {
		f.cells[i][j] = rune(f.input[i][j])
		i += di
		j += dj
	}
}

func (f field) output() (o []string) {
	for _, row := range f.cells {
		o = append(o, string(row))
	}
	return
}

func (d Solve) CorrectAnswers() [2]string {
	return [2]string{"2569", "1998"}
}
