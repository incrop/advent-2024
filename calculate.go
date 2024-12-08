package main

import (
	"incrop/advent-2024/day01"
	"incrop/advent-2024/day02"
	"incrop/advent-2024/day03"
	"incrop/advent-2024/day04"

	tea "github.com/charmbracelet/bubbletea/v2"
)

type Calculate interface {
	Part1(input []string) int64
	Part2(input []string) int64
	Answers() (int64, int64)
}

type AnswerMsg struct {
	day    int
	part   int
	answer int64
}

func calculateCmd(c Calculate, day int, part int, input []string) tea.Cmd {
	return func() tea.Msg {
		answerMsg := AnswerMsg{
			day:    day,
			part:   part,
			answer: 0,
		}
		if part == 0 {
			answerMsg.answer = c.Part1(input)
		} else {
			answerMsg.answer = c.Part2(input)
		}
		return answerMsg
	}
}

func collectCalculations() [26]Calculate {
	return [26]Calculate{
		nil,
		day01.Calculate{},
		day02.Calculate{},
		day03.Calculate{},
		day04.Calculate{},
	}
}

func scheduleAutosolve(calculations [26]Calculate, presets loadedPresets) tea.Cmd {
	var commands []tea.Cmd
	for day, c := range calculations {
		if c == nil {
			continue
		}
		commands = append(
			commands,
			calculateCmd(c, day, 0, presets.input(day, 1)),
			calculateCmd(c, day, 1, presets.input(day, 1)),
		)
	}
	return tea.Batch(commands...)
}
