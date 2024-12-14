package day11

import (
	"fmt"
	"incrop/advent-2024/out"
	"log"
	"strconv"
	"strings"
)

type Calculate struct{}

func (d Calculate) Part1(input []string, outputCh chan<- []string) int64 {
	s := parse(input)
	l := out.NewLog(outputCh)
	for range 25 {
		l.Printf("%s", s)
		s = s.expand()
	}
	l.Printf("%s", s)
	return s.count()
}

func (d Calculate) Part2(input []string, outputCh chan<- []string) int64 {
	s := parse(input)
	l := out.NewLog(outputCh)
	for range 75 {
		l.Printf("%s", s)
		s = s.expand()
	}
	l.Printf("%s", s)
	return s.count()
}

type stones map[int64]int64

func (s1 stones) expand() (s2 stones) {
	s2 = make(stones)
	for stone, count := range s1 {
		if stone == 0 {
			s2[1] += count
			continue
		}
		digits := fmt.Sprintf("%d", stone)
		if len(digits)%2 == 0 {
			next1, err := strconv.ParseInt(digits[:len(digits)/2], 10, 64)
			if err != nil {
				log.Fatal(err)
			}
			next2, err := strconv.ParseInt(digits[len(digits)/2:], 10, 64)
			if err != nil {
				log.Fatal(err)
			}
			s2[next1] += count
			s2[next2] += count
			continue
		}
		s2[stone*2024] += count
	}
	return
}

func (s stones) count() (total int64) {
	for _, count := range s {
		total += count
	}
	return
}

func parse(input []string) (s stones) {
	s = make(stones)
	for _, line := range input {
		for _, field := range strings.Fields(line) {
			num, err := strconv.ParseInt(field, 10, 64)
			if err != nil {
				log.Fatal(err)
			}
			s[num]++
		}
	}
	return
}

func (s stones) String() string {
	var sb strings.Builder
	for stone, count := range s {
		fmt.Fprintf(&sb, "%d:%d ", stone, count)
	}
	return sb.String()
}

func (d Calculate) CorrectAnswers() [2]int64 {
	return [2]int64{189167, 225253278506288}
}
