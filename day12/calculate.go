package day12

type Calculate struct{}

func (d Calculate) Part1(input []string, outputCh chan<- []string) int64 {
	f := parse(input)
	outputCh <- f.output()
	return f.calculateTotalPrice(false)
}

func (d Calculate) Part2(input []string, outputCh chan<- []string) int64 {
	f := parse(input)
	outputCh <- f.output()
	return f.calculateTotalPrice(true)
}

type field []string

func (f field) neighbor(i, j int, off offset) (ii, jj int, rr byte) {
	ii, jj = i+off[0], j+off[1]
	if ii < 0 || ii >= len(f) || jj < 0 || jj >= len(f[ii]) {
		return ii, jj, 0
	}
	return ii, jj, f[ii][jj]
}

type visitedCell struct {
	cell  bool
	walls [4]bool
}

type offset [2]int

var dirOffsets [4]offset = [4]offset{{-1, 0}, {0, 1}, {1, 0}, {0, -1}}

func (f field) calculateTotalPrice(discount bool) (total int64) {
	visited := make([][]visitedCell, len(f))
	for i, line := range f {
		visited[i] = make([]visitedCell, len(line))
	}

	visitWall := func(i, j int, wallDir int) (wallLen int) {
		if visited[i][j].walls[wallDir] {
			return 0
		}
		r := f[i][j]
		wallLen = 1
		visited[i][j].walls[wallDir] = true
		for _, moveDir := range []int{(wallDir + 1) % 4, (wallDir + 3) % 4} {
			moveOff := dirOffsets[moveDir]
			im, jm, rm := i, j, r
			for {
				im, jm, rm = f.neighbor(im, jm, moveOff)
				if rm != r {
					break
				}
				_, _, rw := f.neighbor(im, jm, dirOffsets[wallDir])
				if rw == r {
					break
				}
				visited[im][jm].walls[wallDir] = true
				wallLen++
			}
		}
		return
	}

	var visitRegion func(i, j int) (area, perimeter int)
	visitRegion = func(i, j int) (area, perimeter int) {
		area = 1
		r := f[i][j]
		visited[i][j].cell = true
		for dir, off := range dirOffsets {
			ii, jj, rr := f.neighbor(i, j, off)
			if r != rr {
				wallLen := visitWall(i, j, dir)
				if wallLen > 0 && discount {
					wallLen = 1
				}
				perimeter += wallLen
				continue
			}
			if visited[ii][jj].cell {
				continue
			}
			moreArea, morePerimeter := visitRegion(ii, jj)
			area += moreArea
			perimeter += morePerimeter
		}
		return
	}
	for i, line := range f {
		for j := range line {
			if !visited[i][j].cell {
				area, perimeter := visitRegion(i, j)
				total += int64(area) * int64(perimeter)
			}
		}
	}
	return
}

func (f field) _calculateTotalPrice() (total int64) {
	visited := make([][]bool, len(f))
	for i, line := range f {
		visited[i] = make([]bool, len(line))
	}
	var regionAreaAndPerimeter func(i, j int) (area, perimeter int)
	regionAreaAndPerimeter = func(i, j int) (area, perimeter int) {
		area = 1
		r := f[i][j]
		visited[i][j] = true
		for _, off := range dirOffsets {
			ii, jj := i+off[0], j+off[1]
			if ii < 0 || ii >= len(f) || jj < 0 || jj >= len(f[ii]) {
				perimeter++
				continue
			}
			rr := f[ii][jj]
			if r != rr {
				perimeter++
				continue
			}
			if visited[ii][jj] {
				continue
			}
			moreArea, morePerimeter := regionAreaAndPerimeter(ii, jj)
			area += moreArea
			perimeter += morePerimeter
		}
		return
	}
	for i, line := range f {
		for j := range line {
			if !visited[i][j] {
				area, perimeter := regionAreaAndPerimeter(i, j)
				total += int64(area) * int64(perimeter)
			}
		}
	}
	return
}

func parse(input []string) (f field) {
	for _, line := range input {
		f = append(f, line)
	}
	return
}

func (f field) output() (lines []string) {
	h, w := len(f), len(f[0])
	runes := make([][]rune, h*2+1)
	for i := range runes {
		runes[i] = make([]rune, w*2+1)
		for j := range runes[i] {
			runes[i][j] = ' '
		}
	}
	for i, row := range f {
		for j, r := range row {
			ri, rj := i*2+1, j*2+1
			runes[ri][rj] = r
			for dir, off := range dirOffsets {
				_, _, rr := f.neighbor(i, j, off)
				if rune(rr) != r {
					wall := '|'
					if dir%2 == 0 {
						wall = '-'
					}
					off1, off2 := dirOffsets[(dir+1)%4], dirOffsets[(dir+3)%4]
					runes[ri+off[0]][rj+off[1]] = wall
					runes[ri+off[0]+off1[0]][rj+off[1]+off1[1]] = '+'
					runes[ri+off[0]+off2[0]][rj+off[1]+off2[1]] = '+'
				}
			}
		}
	}
	for _, line := range runes {
		lines = append(lines, string(line))
	}
	return
}

func (d Calculate) CorrectAnswers() [2]int64 {
	return [2]int64{1464678, 877492}
}
