package day20

import (
	"fmt"
	"slices"
)

type Solve struct{}

func (d Solve) Part1(input []string, outputCh chan<- []string) int64 {
	f := parse(input)
	total, output := sumAndSortOutput(f.countCheats(2, 100))
	outputCh <- output
	return total
}

func (d Solve) Part2(input []string, outputCh chan<- []string) int64 {
	f := parse(input)
	total, output := sumAndSortOutput(f.countCheats(20, 100))
	outputCh <- output
	return total
}

func sumAndSortOutput(cheatsBySave map[int]int) (total int64, output []string) {
	total = int64(0)
	var saves []int
	for save := range cheatsBySave {
		saves = append(saves, save)
	}
	slices.Sort(saves)
	for _, save := range saves {
		output = append(output, fmt.Sprintf("save: %d count %d", save, cheatsBySave[save]))
		total += int64(cheatsBySave[save])
	}
	return
}

type coord [2]int

var dirOffsets [4]coord = [4]coord{{-1, 0}, {0, 1}, {1, 0}, {0, -1}}

type field struct {
	cells [][]int
	start coord
	exit  coord
}

func min(a, b int) int {
	if b < a {
		return b
	}
	return a
}

func max(a, b int) int {
	if b > a {
		return b
	}
	return a
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func (f field) countCheats(skip int, minSave int) (cheatsBySave map[int]int) {
	cheatsBySave = map[int]int{}
	iMax, jMax := len(f.cells)-1, len(f.cells[0])-1
	for i1, row := range f.cells {
		if i1 == 0 || i1 == iMax {
			continue
		}
		for j1, v1 := range row {
			if j1 == 0 || j1 == jMax || v1 == -1 {
				continue
			}
			for i2 := max(i1-skip, 1); i2 <= min(i1+skip, iMax-1); i2++ {
				dj := skip - abs(i2-i1)
				for j2 := max(j1-dj, 1); j2 <= min(j1+dj, jMax-1); j2++ {
					dist := abs(i2-i1) + abs(j2-j1)
					if dist < 2 {
						continue
					}
					v2 := f.cells[i2][j2]
					if v2 == -1 {
						continue
					}
					if save := v2 - v1 - dist; save >= minSave {
						cheatsBySave[save]++
					}
				}
			}
		}
	}
	return
}

func parse(input []string) (f field) {
	f.cells = make([][]int, len(input))
	for i, line := range input {
		f.cells[i] = make([]int, len(line))
		for j, r := range line {
			switch r {
			case 'S':
				f.start = [2]int{i, j}
			case 'E':
				f.exit = [2]int{i, j}
			case '#':
				f.cells[i][j] = -1
			}
		}
	}
	curr := f.start
	for n := 1; curr != f.exit; n++ {
		for _, off := range dirOffsets {
			i, j := curr[0]+off[0], curr[1]+off[1]
			next := coord{i, j}
			if f.cells[i][j] != 0 || f.start == next {
				continue
			}
			f.cells[i][j] = n
			curr = next
			break
		}
	}
	return
}

func (d Solve) CorrectAnswers() [2]int64 {
	return [2]int64{1323, 983905}
}
