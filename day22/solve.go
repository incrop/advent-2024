package day22

import (
	"fmt"
	"incrop/advent-2024/out"
	"iter"
	"log"
	"strconv"
)

type Solve struct{}

func (d Solve) Part1(input []string, outputCh chan<- []string) string {
	s := parse(input)
	for range 2000 {
		for i := range s {
			s[i] = nextSecret(s[i])
		}
	}
	outputCh <- s.output()
	sum := int64(0)
	for _, secret := range s {
		sum += secret
	}
	return strconv.FormatInt(sum, 10)
}

func (d Solve) Part2(input []string, outputCh chan<- []string) string {
	s := parse(input)

	profitPerDiffWindow := map[[4]int8]int64{}
	for _, secret := range s {
		firstDiffWindow := map[[4]int8]bool{}
		for price, diffWindow := range priceAndDiffSlidingWindow(secret, 2000) {
			if firstDiffWindow[diffWindow] {
				continue
			}
			firstDiffWindow[diffWindow] = true
			profitPerDiffWindow[diffWindow] += int64(price)
		}
	}
	maxProfit := int64(0)
	maxDiffWindow := [4]int8{}
	for diffWindow, profit := range profitPerDiffWindow {
		if profit > maxProfit {
			maxProfit = profit
			maxDiffWindow = diffWindow
		}
	}
	l := out.NewLog(outputCh)
	for _, secret := range s {
		found := false
		for price, diffWindow := range priceAndDiffSlidingWindow(secret, 2000) {
			if diffWindow == maxDiffWindow {
				l.Printf("%v: price %d", diffWindow, price)
				found = true
				break
			}
		}
		if !found {
			l.Printf("%v: not found", maxDiffWindow)
		}
	}
	return strconv.FormatInt(maxProfit, 10)
}

func priceAndDiffSlidingWindow(secret int64, secretUpdatesCount int) iter.Seq2[int8, [4]int8] {
	return func(yield func(int8, [4]int8) bool) {
		lastPrice := int8(secret % 10)
		diffWindow := [4]int8{}
		for i := range secretUpdatesCount {
			secret = nextSecret(secret)
			price := int8(secret % 10)
			if i < 4 {
				diffWindow[i] = price - lastPrice
			} else {
				diffWindow = [4]int8{diffWindow[1], diffWindow[2], diffWindow[3], int8(price - lastPrice)}
			}
			lastPrice = price
			if i < 3 {
				continue
			}
			if !yield(price, diffWindow) {
				return
			}
		}
	}
}

type secrets []int64

func nextSecret(secret int64) int64 {
	secret = mixAndPrune(secret, secret*64)
	secret = mixAndPrune(secret, secret/32)
	secret = mixAndPrune(secret, secret*2048)
	return secret
}

func mixAndPrune(secret int64, next int64) int64 {
	return (secret ^ next) % 16777216
}

func parse(input []string) (s secrets) {
	for _, line := range input {
		secret, err := strconv.ParseInt(line, 10, 64)
		if err != nil {
			log.Fatal(err)
		}
		s = append(s, secret)
	}
	return
}

func (s secrets) output() (lines []string) {
	for _, secret := range s {
		lines = append(lines, fmt.Sprintf("%d", secret))
	}
	return
}

func (d Solve) CorrectAnswers() [2]string {
	return [2]string{"15303617151", "1727"}
}
