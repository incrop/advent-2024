package day03

import (
	"log"
	"regexp"
	"strconv"
	"strings"
)

type Calculate struct{}

func (d Calculate) Part1(input []string) int64 {
	mulSum := int64(0)
	for _, e := range parse(input) {
		switch e := e.(type) {
		case mul:
			mulSum += int64(e.x * e.y)
		}
	}
	return mulSum
}

func (d Calculate) Part2(input []string) int64 {
	doing := true
	mulSum := int64(0)
	for _, e := range parse(input) {
		switch e := e.(type) {
		case do:
			doing = bool(e)
		case mul:
			if doing {
				mulSum += int64(e.x * e.y)
			}
		}
	}
	return mulSum
}

type expr interface{}

type mul struct {
	x, y int
}

type do bool

var parseRegex = regexp.MustCompile(`mul\((\d{1,3}),(\d{1,3})\)|do\(\)|don't\(\)`)

func parse(input []string) (expr []expr) {
	for _, line := range input {
		for _, match := range parseRegex.FindAllStringSubmatch(line, -1) {
			if strings.HasPrefix(match[0], "do(") {
				expr = append(expr, do(true))
			} else if strings.HasPrefix(match[0], "don't(") {
				expr = append(expr, do(false))
			} else if strings.HasPrefix(match[0], "mul(") {
				x, err := strconv.Atoi(match[1])
				if err != nil {
					log.Fatal(err)
				}
				y, err := strconv.Atoi(match[2])
				if err != nil {
					log.Fatal(err)
				}
				expr = append(expr, mul{x, y})
			}
		}
	}
	return
}

func (d Calculate) Answers() (int64, int64) {
	return 188741603, 67269798
}
