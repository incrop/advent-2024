package day19

import (
	"incrop/advent-2024/out"
	"strings"
)

type Solve struct{}

func (d Solve) Part1(input []string, outputCh chan<- []string) (result int64) {
	r := parse(input)
	l := out.NewLog(outputCh)
	l.Printf("%s", strings.Join(r.towels, ", "))
	l.Printf("")
	for _, design := range r.designs {
		parts := findWayToSplit(design, r.towels)
		if len(parts) > 0 {
			l.Printf("%s", strings.Join(parts, " "))
			result++
		} else {
			l.Printf("impossible")
		}
	}
	return
}

func findWayToSplit(design string, towels []string) []string {
	cache := map[string][]string{}
	var recurse func(design string) []string
	recurse = func(design string) (result []string) {
		if cached, ok := cache[design]; ok {
			return cached
		}
		defer func() {
			cache[design] = result
		}()
		for _, towel := range towels {
			if design == towel {
				return []string{design}
			}
			if !strings.HasPrefix(design, towel) {
				continue
			}
			parts := recurse(strings.TrimPrefix(design, towel))
			if len(parts) > 0 {
				return append([]string{towel}, parts...)
			}
		}
		return nil
	}
	return recurse(design)
}

func (d Solve) Part2(input []string, outputCh chan<- []string) (result int64) {
	r := parse(input)
	l := out.NewLog(outputCh)
	l.Printf("%s", strings.Join(r.towels, ", "))
	l.Printf("")
	for _, design := range r.designs {
		waysToSplit := countWaysToSplit(design, r.towels)
		l.Printf("ways to split: %d", waysToSplit)
		result += waysToSplit
	}
	return
}

func countWaysToSplit(design string, towels []string) int64 {
	cache := map[string]int64{}
	var recurse func(design string) int64
	recurse = func(design string) int64 {
		if cached, ok := cache[design]; ok {
			return cached
		}
		result := int64(0)
		for _, towel := range towels {
			if design == towel {
				result += 1
				continue
			}
			if !strings.HasPrefix(design, towel) {
				continue
			}
			result += recurse(strings.TrimPrefix(design, towel))
		}
		cache[design] = result
		return result
	}
	return recurse(design)
}

type requirements struct {
	towels  []string
	designs []string
}

func parse(input []string) (r requirements) {
	r.towels = strings.Split(input[0], ", ")
	r.designs = input[2:]
	return
}

func (d Solve) CorrectAnswers() [2]int64 {
	return [2]int64{206, 622121814629343}
}
