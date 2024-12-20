package day05

import (
	"incrop/advent-2024/out"
	"log"
	"strconv"
	"strings"
)

type Solve struct{}

type ordering struct {
	before int
	after  int
}

type update []int

func (u update) isValid(bannedOrderings map[ordering]bool) bool {
	for i, before := range u[:len(u)-1] {
		for _, after := range u[i+1:] {
			if bannedOrderings[ordering{before, after}] {
				return false
			}
		}
	}
	return true
}

func (u update) fixInvalid(bannedOrderings map[ordering]bool) (invalidFound bool) {
	for i := range u[:len(u)-1] {
	retryIteration:
		for {
			before := u[i]
			for j := i + 1; j < len(u); j++ {
				after := u[j]
				if bannedOrderings[ordering{before, after}] {
					u[i], u[j] = u[j], u[i]
					invalidFound = true
					continue retryIteration
				}
			}
			break
		}
	}
	return
}

type rules struct {
	orderings []ordering
	updates   []update
}

func (d Solve) Part1(input []string, outputCh chan<- []string) int64 {
	rules := parse(input)
	middleNumSum := int64(0)
	bannedOrderings := make(map[ordering]bool)
	for _, o := range rules.orderings {
		bannedOrderings[ordering{o.after, o.before}] = true
	}
	l := out.NewLog(outputCh)
	for _, update := range rules.updates {
		if update.isValid(bannedOrderings) {
			middleNum := update[(len(update)-1)/2]
			middleNumSum += int64(middleNum)
			l.Printf("valid, middle: %d", middleNum)
		} else {
			l.Printf("invalid")
		}
	}
	return middleNumSum
}

func (d Solve) Part2(input []string, outputCh chan<- []string) int64 {
	rules := parse(input)
	middleNumSum := int64(0)
	bannedOrderings := make(map[ordering]bool)
	for _, o := range rules.orderings {
		bannedOrderings[ordering{o.after, o.before}] = true
	}
	l := out.NewLog(outputCh)
	for _, update := range rules.updates {
		if update.fixInvalid(bannedOrderings) {
			middleNum := update[(len(update)-1)/2]
			middleNumSum += int64(middleNum)
			l.Printf("invalid, middle after fix: %d", middleNum)
		} else {
			l.Printf("valid, skipping")
		}
	}
	return middleNumSum
}

func parse(input []string) (rules rules) {
	var nextIdx int
	rules.orderings, nextIdx = parseOrderings(input)
	rules.updates = parseUpdates(input[nextIdx:])
	return
}

func parseOrderings(input []string) (orderings []ordering, nextIdx int) {
	for i, line := range input {
		if line == "" {
			return orderings, i + 1
		}
		beforeAndAfter := strings.Split(line, "|")
		before, err := strconv.Atoi(beforeAndAfter[0])
		if err != nil {
			log.Fatal(err)
		}
		after, err := strconv.Atoi(beforeAndAfter[1])
		if err != nil {
			log.Fatal(err)
		}
		orderings = append(orderings, ordering{
			before: before,
			after:  after,
		})
	}
	return orderings, len(input)
}

func parseUpdates(input []string) (updates []update) {
	for _, line := range input {
		pagesStr := strings.Split(line, ",")
		update := make(update, len(pagesStr))
		for i, pageStr := range pagesStr {
			page, err := strconv.Atoi(pageStr)
			if err != nil {
				log.Fatal(err)
			}
			update[i] = page
		}
		updates = append(updates, update)
	}
	return
}

func (d Solve) CorrectAnswers() [2]int64 {
	return [2]int64{5452, 4598}
}
