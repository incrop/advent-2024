package day01

import (
	"incrop/advent-2024/out"
	"log"
	"slices"
	"strconv"
	"strings"
)

type Calculate struct{}

func (d Calculate) Part1(input []string, outputCh chan<- []string) int64 {
	left, right := parse(input)
	slices.Sort(left)
	slices.Sort(right)
	diffSum := int64(0)
	l := out.NewLog(outputCh)
	for i, numLeft := range left {
		numRight := right[i]
		diff := numLeft - numRight
		if diff < 0 {
			diff = -diff
		}
		diffSum += int64(diff)
		l.Printf("abs(%d - %d) = %d", numLeft, numRight, diff)
	}
	return diffSum
}

func (d Calculate) Part2(input []string, outputCh chan<- []string) int64 {
	left, right := parse(input)
	rightFreq := make(map[int]int)
	for _, numRight := range right {
		rightFreq[numRight] += 1
	}
	similarityScore := int64(0)
	l := out.NewLog(outputCh)
	for _, numLeft := range left {
		similarity := numLeft * rightFreq[numLeft]
		similarityScore += int64(similarity)
		l.Printf("%d * %d = %d", numLeft, rightFreq[numLeft], similarity)
	}
	return similarityScore
}

func parse(input []string) ([]int, []int) {
	var left, right []int
	for _, line := range input {
		fields := strings.Fields(line)
		numLeft, err := strconv.Atoi(fields[0])
		if err != nil {
			log.Fatal(err)
		}
		numRight, err := strconv.Atoi(fields[1])
		if err != nil {
			log.Fatal(err)
		}
		left = append(left, numLeft)
		right = append(right, numRight)
	}
	return left, right
}

func (d Calculate) CorrectAnswers() [2]int64 {
	return [2]int64{765748, 27732508}
}
