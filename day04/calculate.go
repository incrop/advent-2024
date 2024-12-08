package day04

type Calculate struct{}

func (d Calculate) Part1(input []string) int64 {
	xmasCount := int64(0)
	for i, line := range input {
		for j := range line {
			xmasCount += int64(xmasCountStaringAt(input, i, j))
		}
	}
	return xmasCount
}

func xmasCountStaringAt(input []string, i, j int) (count int) {
	for di := -1; di <= 1; di++ {
		for dj := -1; dj <= 1; dj++ {
			if di == 0 && dj == 0 {
				continue
			}
			if readsXmsInThisDirection(input, i, j, di, dj) {
				count++
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

func (d Calculate) Part2(input []string) int64 {
	xmasCount := int64(0)
	for i, line := range input {
		for j := range line {
			if hasMasInXShapeAt(input, i, j) {
				xmasCount++
			}
		}
	}
	return xmasCount
}

func hasMasInXShapeAt(input []string, i, j int) bool {
	if input[i][j] != 'A' {
		return false
	}
	for di, dj := -1, -1; di <= 1; di += 2 {
		msCount := [2]int{}
		for _, coord := range [2][2]int{{i+di, j+dj}, {i-di, j-dj}} {
			if coord[0] < 0 || coord[0] >= len(input) {
				return false
			}
			line := input[coord[0]]
			if coord[1] < 0 || coord[1] >= len(line) {
				return false
			}
			switch line[coord[1]]{
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

func (d Calculate) Answers() (int64, int64) {
	return 2569, 1998
}
