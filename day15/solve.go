package day15

import (
	"strconv"
	"strings"
)

type Solve struct{}

func (d Solve) Part1(input []string, outputCh chan<- []string) string {
	f := parse(input, false)
	outputCh <- f.output()
	i := 0
	for f.executeNextCommand() {
		i++
		if i%8 == 0 {
			outputCh <- f.output()
		}
	}
	outputCh <- f.output()
	return strconv.FormatInt(f.boxesGpsCoordSum(), 10)
}

func (d Solve) Part2(input []string, outputCh chan<- []string) string {
	f := parse(input, true)
	outputCh <- f.output()
	i := 0
	for f.executeNextCommand() {
		i++
		if i%8 == 0 {
			outputCh <- f.output()
		}
	}
	return strconv.FormatInt(f.boxesGpsCoordSum(), 10)
}

type coord [2]int

func (c1 coord) plus(c2 coord) coord {
	return coord{c1[0] + c2[0], c1[1] + c2[1]}
}

type field struct {
	cells    [][]rune
	pos      coord
	nextCmd  int
	commands []command
}

type command int

var commandRunes = "^>v<"
var commandOffsets [4]coord = [4]coord{{-1, 0}, {0, 1}, {1, 0}, {0, -1}}

func (f *field) cell(c coord) *rune {
	return &f.cells[c[0]][c[1]]
}

func (f *field) executeNextCommand() bool {
	if f.nextCmd >= len(f.commands) {
		return false
	}
	cmd := f.commands[f.nextCmd]
	off := commandOffsets[cmd]
	f.tryMove(off)
	f.nextCmd++
	return true
}

func (f *field) tryMove(off coord) {
	acc := []coord{}
	next := map[coord]bool{f.pos: true}
	for len(next) > 0 {
		for pos := range next {
			acc = append(acc, pos)
		}
		curr := next
		next = map[coord]bool{}
		for pos := range curr {
			nextPos := pos.plus(off)
			r := *f.cell(nextPos)
			if r == '#' {
				return
			}
			if r == '.' {
				continue
			}
			next[nextPos] = true
			if off[0] == 0 {
				continue
			}
			if r == '[' {
				next[nextPos.plus(coord{0, 1})] = true
			} else if r == ']' {
				next[nextPos.plus(coord{0, -1})] = true
			}
		}
	}
	for i := len(acc) - 1; i >= 0; i-- {
		p1 := acc[i]
		p2 := p1.plus(off)
		c1, c2 := f.cell(p1), f.cell(p2)
		*c1, *c2 = *c2, *c1
	}
	f.pos = f.pos.plus(off)
}

func (f *field) boxesGpsCoordSum() (sum int64) {
	for i, row := range f.cells {
		for j, r := range row {
			if r == 'O' || r == '[' {
				sum += int64(100*i + j)
			}
		}
	}
	return
}

func parse(input []string, expand bool) (f field) {
	nextIndex := 0
	if expand {
		nextIndex = f.parseExpandedCells(input)
	} else {
		nextIndex = f.parseCells(input)
	}
	f.parseCommands(input[nextIndex:])
	return
}

func (f *field) parseCells(input []string) (nextIndex int) {
	for i, line := range input {
		if line == "" {
			return i + 1
		}
		row := make([]rune, len(line))
		for j, r := range line {
			if r == '@' {
				f.pos = coord{i, j}
			}
			row[j] = r
		}
		f.cells = append(f.cells, row)
	}
	return len(input)
}

func (f *field) parseExpandedCells(input []string) (nextIndex int) {
	for i, line := range input {
		if line == "" {
			return i + 1
		}
		row := make([]rune, len(line)*2)
		for j, r := range line {
			j1, j2 := j*2, j*2+1
			switch r {
			case '@':
				f.pos = coord{i, j * 2}
				row[j1], row[j2] = '@', '.'
			case '.':
				row[j1], row[j2] = '.', '.'
			case '#':
				row[j1], row[j2] = '#', '#'
			case 'O':
				row[j1], row[j2] = '[', ']'
			}
		}
		f.cells = append(f.cells, row)
	}
	return len(input)
}

func (f *field) parseCommands(input []string) {
	for _, line := range input {
		for _, r := range line {
			cmd := command(strings.IndexRune(commandRunes, r))
			f.commands = append(f.commands, cmd)
		}
	}
}

func (f field) output() (lines []string) {
	for _, row := range f.cells {
		lines = append(lines, string(row))
	}
	cmdlineWidth := 64
	cursorPos := f.nextCmd % cmdlineWidth
	lines = append(lines, strings.Repeat(" ", cursorPos)+"↓")
	for i := 0; i < len(f.commands); i += cmdlineWidth {
		j := i + cmdlineWidth
		if j > len(f.commands) {
			j = len(f.commands)
		}
		if j < f.nextCmd {
			continue
		}
		cmdline := make([]byte, j-i)
		for k, cmd := range f.commands[i:j] {
			cmdline[k] = commandRunes[cmd]
		}
		lines = append(lines, string(cmdline))
	}
	return
}

func (d Solve) CorrectAnswers() [2]string {
	return [2]string{"1442192", "1448458"}
}
