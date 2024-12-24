package day13

import (
	"incrop/advent-2024/out"
	"log"
	"regexp"
	"strconv"
)

type Solve struct{}

func (d Solve) Part1(input []string, outputCh chan<- []string) string {
	ms := parse(input)
	l := out.NewLog(outputCh)
	return strconv.FormatInt(ms.solutionsCostSum(l), 10)
}

func (d Solve) Part2(input []string, outputCh chan<- []string) string {
	ms := parse(input)
	for i := range ms {
		ms[i].prize[0] += 10000000000000
		ms[i].prize[1] += 10000000000000
	}
	l := out.NewLog(outputCh)
	return strconv.FormatInt(ms.solutionsCostSum(l), 10)
}

type machines []machine

type coord [2]int64

type machine struct {
	offA  coord
	offB  coord
	prize coord
}

func (ms machines) solutionsCostSum(l *out.Log) (totalCost int64) {
	for _, m := range ms {
		if pressA, pressB, found := m.findSolution(); found {
			cost := pressA*3 + pressB
			l.Printf("%d * [%d %d] = [%d %d]", pressA, m.offA[0], m.offA[1], m.prize[0], m.prize[1])
			l.Printf("%d * [%d %d] = [%d %d]", pressB, m.offB[0], m.offB[1], m.prize[0], m.prize[1])
			l.Printf("%d * 3 + %d = %d", pressA, pressB, cost)
			l.Printf("")
			totalCost += cost
		} else {
			l.Printf("No solution")
			l.Printf("")
			l.Printf("")
			l.Printf("")
		}
	}
	return
}

// Solve for whole and positive pressA and pressB:
// pressA*m.offA[0]+pressB*m.offB[0] == m.prize[0]
// pressA*m.offA[1]+pressB*m.offB[1] == m.prize[1]
func (m machine) findSolution() (pressA, pressB int64, found bool) {
	topA := m.prize[0]*m.offB[1] - m.prize[1]*m.offB[0]
	botA := m.offA[0]*m.offB[1] - m.offA[1]*m.offB[0]
	if topA%botA != 0 {
		return 0, 0, false
	}
	pressA = topA / botA
	if pressA < 0 {
		return 0, 0, false
	}
	topB := m.prize[0] - pressA*m.offA[0]
	botB := m.offB[0]
	if topB%botB != 0 {
		return 0, 0, false
	}
	pressB = topB / botB
	if pressB < 0 {
		return 0, 0, false
	}
	return pressA, pressB, true
}

var buttonRegexp = regexp.MustCompile(`^Button ([A|B]): X\+(\d+), Y\+(\d+)$`)
var prizeRegexp = regexp.MustCompile(`^Prize: X=(\d+), Y=(\d+)$`)

func parse(input []string) (ms machines) {
	var m machine
	for _, line := range input {
		if line == "" {
			continue
		}
		if match := buttonRegexp.FindStringSubmatch(line); match != nil {
			dx, err := strconv.ParseInt(match[2], 10, 64)
			if err != nil {
				log.Fatal(err)
			}
			dy, err := strconv.ParseInt(match[3], 10, 64)
			if err != nil {
				log.Fatal(err)
			}
			switch match[1] {
			case "A":
				m.offA = coord{dx, dy}
			case "B":
				m.offB = coord{dx, dy}
			}
			continue
		}
		if match := prizeRegexp.FindStringSubmatch(line); match != nil {
			x, err := strconv.ParseInt(match[1], 10, 64)
			if err != nil {
				log.Fatal(err)
			}
			y, err := strconv.ParseInt(match[2], 10, 64)
			if err != nil {
				log.Fatal(err)
			}
			m.prize = coord{x, y}
			ms = append(ms, m)
			m = machine{}
			continue
		}
	}
	return
}

func (d Solve) CorrectAnswers() [2]string {
	return [2]string{"36571", "85527711500010"}
}
