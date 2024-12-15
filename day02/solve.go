package day02

import (
	"incrop/advent-2024/out"
	"log"
	"strconv"
	"strings"
)

type Solve struct{}

type report []int

func (d Solve) Part1(input []string, outputCh chan<- []string) int64 {
	safeReports := 0
	l := out.NewLog(outputCh)
	for _, r := range parse(input) {
		unsafeIdx := r.unsafeIdx()
		if unsafeIdx == 0 {
			safeReports++
			l.Printf("safe")
		} else {
			l.Printf("unsafe at idx %d", unsafeIdx)
		}
	}
	return int64(safeReports)
}

func (d Solve) Part2(input []string, outputCh chan<- []string) int64 {
	safeReports := 0
	l := out.NewLog(outputCh)
	for _, r := range parse(input) {
		unsafeIdx := r.unsafeIdx()
		if unsafeIdx == 0 {
			l.Printf("safe")
			safeReports++
			continue
		}
		var dampenedTry1 report
		dampenedTry1 = append(dampenedTry1, r[0:unsafeIdx-1]...)
		dampenedTry1 = append(dampenedTry1, r[unsafeIdx:]...)
		if dampenedTry1.unsafeIdx() == 0 {
			l.Printf("safe after removing %d at idx %d", r[unsafeIdx-1], unsafeIdx-1)
			safeReports++
			continue
		}
		var dampenedTry2 report
		dampenedTry2 = append(dampenedTry2, r[0:unsafeIdx]...)
		dampenedTry2 = append(dampenedTry2, r[unsafeIdx+1:]...)
		if dampenedTry2.unsafeIdx() == 0 {
			l.Printf("safe after removing %d at idx %d", r[unsafeIdx], unsafeIdx)
			safeReports++
			continue
		}
		if unsafeIdx != 2 {
			l.Printf("unsafe at idx %d", unsafeIdx)
			continue
		}
		if r[1:].unsafeIdx() == 0 {
			l.Printf("safe after removing %d at idx %d", r[0], 0)
			safeReports++
			continue
		}
	}
	return int64(safeReports)
}

func (report report) unsafeIdx() int {
	if len(report) < 2 || report[0] == report[1] {
		return 1
	}
	diff := report[1] - report[0]
	if diff < -3 || diff > 3 {
		return 1
	}
	isIncr := report[0] < report[1]
	prev := report[1]
	for i, next := range report[2:] {
		if prev == next || (prev < next != isIncr) {
			return i + 2
		}
		diff := next - prev
		if diff < -3 || diff > 3 {
			return i + 2
		}
		prev = next
	}
	return 0
}

func parse(input []string) (reports []report) {
	for _, line := range input {
		fields := strings.Fields(line)
		var report report
		for _, field := range fields {
			level, err := strconv.Atoi(field)
			if err != nil {
				log.Fatal(err)
			}
			report = append(report, level)
		}
		reports = append(reports, report)
	}
	return
}

func (d Solve) CorrectAnswers() [2]int64 {
	return [2]int64{236, 308}
}
