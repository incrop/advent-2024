package day18

import (
	"fmt"
	"log"
	"slices"
	"strconv"
	"strings"
)

type Solve struct{}

func (d Solve) Part1(input []string, outputCh chan<- []string) string {
	f := parse(input)
	path := f.findShortestPath(f.firstBytes)
	outputCh <- f.output(f.firstBytes, path)
	return strconv.Itoa(len(path) - 1)
}

func (d Solve) Part2(input []string, outputCh chan<- []string) string {
	f := parse(input)
	outputCh <- f.output(f.firstBytes, nil)
	minBlocked, maxBlocked := 0, len(f.byteCoords)
	for maxBlocked-minBlocked > 1 {
		mid := (minBlocked + maxBlocked) / 2
		path := f.findShortestPath(mid)
		outputCh <- f.output(mid, path)
		if len(path) > 0 {
			minBlocked = mid
		} else {
			maxBlocked = mid
		}
	}
	path := f.findShortestPath(minBlocked)
	outputCh <- f.output(minBlocked, path)
	blockCoord := f.byteCoords[minBlocked]
	return fmt.Sprintf("%d,%d", blockCoord[0], blockCoord[1])
}

type coord [2]int

type field struct {
	size       coord
	byteCoords []coord
	firstBytes int
}

var dirOffsets [4]coord = [4]coord{{-1, 0}, {0, 1}, {1, 0}, {0, -1}}

func (f *field) findShortestPath(blockedCount int) (path []coord) {
	target := coord{f.size[0] - 1, f.size[1] - 1}
	blocked := map[coord]bool{}
	for _, coord := range f.byteCoords[:blockedCount] {
		blocked[coord] = true
	}
	visited := map[coord]coord{}
	curr := map[coord]coord{
		{0, 0}: {0, 0},
	}
mainLoop:
	for {
		next := map[coord]coord{}
		for currCoord, prevCoord := range curr {
			visited[currCoord] = prevCoord
		}
		for currCoord := range curr {
			for _, off := range dirOffsets {
				nextCoord := coord{currCoord[0] + off[0], currCoord[1] + off[1]}
				if nextCoord[0] < 0 || nextCoord[0] >= f.size[0] || nextCoord[1] < 0 || nextCoord[1] >= f.size[1] {
					continue
				}
				if blocked[nextCoord] {
					continue
				}
				if _, ok := visited[nextCoord]; ok {
					continue
				}
				if nextCoord == target {
					visited[target] = currCoord
					break mainLoop
				}
				next[nextCoord] = currCoord
			}
		}
		if len(next) == 0 {
			return nil
		}
		curr = next
	}
	if _, ok := visited[target]; !ok {
		return
	}
	currCoord := target
	for {
		path = append(path, currCoord)
		prevCoord, ok := visited[currCoord]
		if !ok || prevCoord == currCoord {
			break
		}
		currCoord = prevCoord
	}
	slices.Reverse(path)
	return
}

func parse(input []string) (f *field) {
	f = new(field)
	f.size = parseCoord(input[0])
	firstBytes, err := strconv.Atoi(input[1])
	if err != nil {
		log.Fatal(err)
	}
	f.firstBytes = firstBytes
	for _, line := range input[2:] {
		f.byteCoords = append(f.byteCoords, parseCoord(line))
	}
	return
}

func parseCoord(line string) (c coord) {
	for i, numStr := range strings.Split(line, ",") {
		num, err := strconv.Atoi(numStr)
		if err != nil {
			log.Fatal(err)
		}
		c[i] = num
	}
	return
}

func (f *field) output(numBlocked int, path []coord) (lines []string) {
	cells := make([][]rune, f.size[0])
	for i := range cells {
		cells[i] = make([]rune, f.size[1])
		for j := range cells[i] {
			cells[i][j] = '.'
		}
	}
	for _, coord := range f.byteCoords[:numBlocked] {
		cells[coord[0]][coord[1]] = '#'
	}
	for _, coord := range path {
		cells[coord[0]][coord[1]] = 'O'
	}
	for _, row := range cells {
		lines = append(lines, string(row))
	}
	return
}

func (d Solve) CorrectAnswers() [2]string {
	return [2]string{"252", "5,60"}
}
