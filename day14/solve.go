package day14

import (
	"log"
	"regexp"
	"strconv"
	"strings"
)

type Solve struct{}

func (d Solve) Part1(input []string, outputCh chan<- []string) string {
	f := parse(input)
	for range 100 {
		outputCh <- f.output()
		f.progressOneSecond()
	}
	outputCh <- f.output()
	return strconv.FormatInt(f.safetyFactor(), 10)
}

func (d Solve) Part2(input []string, outputCh chan<- []string) string {
	f := parse(input)
	outputCh <- f.output()
	bestCohesiveness := f.cohesiveness()
	secondsElapsed := int64(0)
	for secondsElapsed < int64(f.h*f.w) {
		f.progressOneSecond()
		secondsElapsed++
		if cohesiveness := f.cohesiveness(); cohesiveness > bestCohesiveness {
			outputCh <- f.output()
			if cohesiveness > bestCohesiveness*2 {
				break
			}
			bestCohesiveness = cohesiveness
		}
	}
	return strconv.FormatInt(secondsElapsed, 10)
}

type robot struct {
	x, y, dx, dy int
}

type field struct {
	w, h   int
	robots []robot
}

func (f field) progressOneSecond() {
	for i := range f.robots {
		r := &f.robots[i]
		r.x = modulo(r.x+r.dx, f.w)
		r.y = modulo(r.y+r.dy, f.h)
	}
}

func modulo(a, b int) int {
	return (a%b + b) % b
}

func (f field) safetyFactor() int64 {
	cx, cy := f.w/2, f.h/2
	quadCounts := [4]int64{}
	for _, r := range f.robots {
		switch [4]bool{r.x < cx, r.x > cx, r.y < cy, r.y > cy} {
		case [4]bool{true, false, true, false}:
			quadCounts[0]++
		case [4]bool{false, true, true, false}:
			quadCounts[1]++
		case [4]bool{true, false, false, true}:
			quadCounts[2]++
		case [4]bool{false, true, false, true}:
			quadCounts[3]++
		}
	}
	return quadCounts[0] * quadCounts[1] * quadCounts[2] * quadCounts[3]
}

func (f field) cohesiveness() (pl int) {
	counts := make(map[[2]int]int)
	for _, r := range f.robots {
		counts[[2]int{r.x, r.y}]++
	}
	for _, r := range f.robots {
		for _, off := range [4][2]int{{-1, 0}, {0, 1}, {1, 0}, {0, -1}} {
			pl += counts[[2]int{r.x + off[0], r.y + off[1]}]
		}
	}
	return
}

var robotRegexp = regexp.MustCompile(`^p=(-?\d+),(-?\d+) v=(-?\d+),(-?\d+)`)

func parse(input []string) (f field) {
	widthAndHeight := strings.Fields(input[0])
	width, err := strconv.Atoi(widthAndHeight[0])
	if err != nil {
		log.Fatal(err)
	}
	height, err := strconv.Atoi(widthAndHeight[1])
	if err != nil {
		log.Fatal(err)
	}
	f.w, f.h = width, height
	for _, line := range input[1:] {
		match := robotRegexp.FindStringSubmatch(line)
		coords := [4]int{}
		for i := range coords {
			coord, err := strconv.Atoi(match[i+1])
			if err != nil {
				log.Fatal(err)
			}
			coords[i] = coord
		}
		f.robots = append(f.robots, robot{
			x:  coords[0],
			y:  coords[1],
			dx: coords[2],
			dy: coords[3],
		})
	}
	return
}

func (f field) output() (lines []string) {
	counts := make([][]int, f.h)
	for y := range counts {
		counts[y] = make([]int, f.w)
	}
	for _, r := range f.robots {
		counts[r.y][r.x]++
	}
	cx, cy := f.w/2, f.h/2
	for y, row := range counts {
		var sb strings.Builder
		for x, count := range row {
			if x == cx || y == cy {
				sb.WriteRune(' ')
				continue
			}
			if count == 0 {
				sb.WriteRune('.')
				continue
			}
			lastDigit := strconv.Itoa(count % 10)
			sb.WriteString(lastDigit)
		}
		lines = append(lines, sb.String())
	}
	return lines
}

func (d Solve) CorrectAnswers() [2]string {
	return [2]string{"220971520", "6355"}
}
