package main

import (
	"incrop/advent-2024/day01"
	"incrop/advent-2024/day02"
	"incrop/advent-2024/day03"
	"incrop/advent-2024/day04"
	"incrop/advent-2024/day05"

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

func (d dayState) calculateCmd(preset int, part int) tea.Cmd {
	return func() tea.Msg {
		answerMsg := AnswerMsg{
			day:    d.day,
			part:   part,
			answer: 0,
		}
		input := d.presets.input(preset)
		if part == 0 {
			answerMsg.answer = d.calculate.Part1(input)
		} else {
			answerMsg.answer = d.calculate.Part2(input)
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
		day05.Calculate{},
	}
}

func (m model) scheduleAutosolve() tea.Cmd {
	var commands []tea.Cmd
	for _, d := range m.dayStates {
		if d.calculate == nil {
			continue
		}
		commands = append(
			commands,
			d.calculateCmd(1, 0),
			d.calculateCmd(1, 1),
		)
	}
	return tea.Batch(commands...)
}
