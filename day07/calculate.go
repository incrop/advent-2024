package day07

import (
	"fmt"
	"incrop/advent-2024/out"
	"log"
	"strconv"
	"strings"
)

type Calculate struct{}

func (d Calculate) Part1(input []string, outputCh chan<- []string) int64 {
	calibrationResult := int64(0)
	l := out.NewLog(outputCh)
	for _, equation := range parse(input) {
		if ops := equation.solution(false); ops != nil {
			l.Printf("%s", equation.print(ops))
			calibrationResult += equation.result
		}
	}
	return calibrationResult
}

func (d Calculate) Part2(input []string, outputCh chan<- []string) int64 {
	calibrationResult := int64(0)
	l := out.NewLog(outputCh)
	for _, equation := range parse(input) {
		if ops := equation.solution(true); ops != nil {
			l.Printf("%s", equation.print(ops))
			calibrationResult += equation.result
		}
	}
	return calibrationResult
}

type equation struct {
	result  int64
	numbers []int64
}

func (eq equation) solution(useConcat bool) []string {
	ops := make([]string, len(eq.numbers)-1)
	var checkRec func(acc int64, i int) bool
	checkRec = func(acc int64, i int) bool {
		if i == len(eq.numbers) {
			return acc == eq.result
		}
		if acc > eq.result {
			return false
		}
		if checkRec(acc+eq.numbers[i], i+1) {
			ops[i-1] = "+"
			return true
		}
		if checkRec(acc*eq.numbers[i], i+1) {
			ops[i-1] = "*"
			return true
		}
		if useConcat && checkRec(concat(acc, eq.numbers[i]), i+1) {
			ops[i-1] = "||"
			return true
		}
		return false
	}
	if checkRec(eq.numbers[0], 1) {
		return ops
	}
	return nil
}

func concat(a, b int64) int64 {
	s := fmt.Sprintf("%d%d", a, b)
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		log.Fatal(err)
	}
	return i
}

func (eq equation) print(ops []string) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d = %d", eq.result, eq.numbers[0])
	for i, number := range eq.numbers[1:] {
		fmt.Fprintf(&sb, " %s %d", ops[i], number)
	}
	return sb.String()
}

func parse(input []string) (eqs []equation) {
	for _, line := range input {
		fields := strings.Fields(line)
		result, err := strconv.Atoi(strings.TrimSuffix(fields[0], ":"))
		if err != nil {
			log.Fatal(err)
		}
		numbers := make([]int64, len(fields)-1)
		for i, field := range fields[1:] {
			number, err := strconv.Atoi(field)
			if err != nil {
				log.Fatal(err)
			}
			numbers[i] = int64(number)
		}
		eqs = append(eqs, equation{
			result:  int64(result),
			numbers: numbers,
		})
	}
	return
}

func (d Calculate) CorrectAnswers() [2]int64 {
	return [2]int64{4364915411363, 38322057216320}
}
