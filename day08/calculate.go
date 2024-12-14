package day08

import "iter"

type Calculate struct{}

func (d Calculate) Part1(input []string, outputCh chan<- []string) int64 {
	f := parse(input)
	antinodeCoords := make(map[coord]bool)
	for _, coords := range f.anntennaCoords {
		for c1, c2 := range coords.allCombinations() {
			diff := c2.minus(c1)
			for _, anc := range [2]coord{c2.plus(diff), c1.minus(diff)} {
				if f.isInside(anc) {
					antinodeCoords[anc] = true
				}
			}
			outputCh <- f.output(antinodeCoords)
		}
	}
	return int64(len(antinodeCoords))
}

func (d Calculate) Part2(input []string, outputCh chan<- []string) int64 {
	f := parse(input)
	antinodeCoords := make(map[coord]bool)
	for _, coords := range f.anntennaCoords {
		for c1, c2 := range coords.allCombinations() {
			diff := c2.minus(c1)
			for anc := c2; f.isInside(anc); anc = anc.plus(diff) {
				antinodeCoords[anc] = true
			}
			for anc := c1; f.isInside(anc); anc = anc.minus(diff) {
				antinodeCoords[anc] = true
			}
			outputCh <- f.output(antinodeCoords)
		}
	}
	return int64(len(antinodeCoords))
}

type coord [2]int

func (c1 coord) plus(c2 coord) coord {
	return coord{c1[0] + c2[0], c1[1] + c2[1]}
}

func (c1 coord) minus(c2 coord) coord {
	return coord{c1[0] - c2[0], c1[1] - c2[1]}
}

type coords []coord

func (coords coords) allCombinations() iter.Seq2[coord, coord] {
	return func(yield func(c1, c2 coord) bool) {
		if len(coords) < 2 {
			return
		}
		for i, c1 := range coords[:len(coords)-1] {
			for _, c2 := range coords[i+1:] {
				if !yield(c1, c2) {
					return
				}
			}
		}
	}
}

type field struct {
	height, width  int
	anntennaCoords map[rune]coords
}

func parse(input []string) (f field) {
	f.height = len(input)
	f.width = len(input[0])
	f.anntennaCoords = make(map[rune]coords)
	for i, line := range input {
		for j, r := range line {
			if r == '.' {
				continue
			}
			f.anntennaCoords[r] = append(f.anntennaCoords[r], coord{i, j})
		}
	}
	return
}

func (f field) isInside(c coord) bool {
	if c[0] < 0 || c[0] >= f.height {
		return false
	}
	if c[1] < 0 || c[1] >= f.width {
		return false
	}
	return true
}

func (f field) output(antinodeCoords map[coord]bool) (lines []string) {
	o := make([][]rune, f.height)
	for i := range o {
		row := make([]rune, f.width)
		for j := range row {
			row[j] = '.'
		}
		o[i] = row
	}
	for coord := range antinodeCoords {
		o[coord[0]][coord[1]] = '#'
	}
	for a, coords := range f.anntennaCoords {
		for _, coord := range coords {
			o[coord[0]][coord[1]] = a
		}
	}
	for _, row := range o {
		lines = append(lines, string(row))
	}
	return
}

func (d Calculate) CorrectAnswers() [2]int64 {
	return [2]int64{318, 1126}
}
