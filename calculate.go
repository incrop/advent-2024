package main

import (
	"incrop/advent-2024/day01"
	"incrop/advent-2024/day02"
	"incrop/advent-2024/day03"

	tea "github.com/charmbracelet/bubbletea/v2"
)

type Calculate interface {
	Part1(input []string) int64
	Part2(input []string) int64
}

type AnswerMsg int64

func calculateCmd(input []string, c Calculate, part int) tea.Cmd {
	return func() tea.Msg {
		if part == 0 {
			return AnswerMsg(c.Part1(input))
		} else {
			return AnswerMsg(c.Part2(input))
		}
	}
}

func collectCalculations() [25]Calculate {
	return [25]Calculate{
		day01.Calculate{},
		day02.Calculate{},
		day03.Calculate{},
	}
}
