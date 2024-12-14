package main

import (
	"incrop/advent-2024/day01"
	"incrop/advent-2024/day02"
	"incrop/advent-2024/day03"
	"incrop/advent-2024/day04"
	"incrop/advent-2024/day05"
	"incrop/advent-2024/day06"
	"incrop/advent-2024/day07"
	"incrop/advent-2024/day08"
	"incrop/advent-2024/day09"
	"incrop/advent-2024/day10"
	"incrop/advent-2024/day11"

	tea "github.com/charmbracelet/bubbletea/v2"
)

type Calculate interface {
	Part1(input []string, outputCh chan<- []string) int64
	Part2(input []string, outputCh chan<- []string) int64
	CorrectAnswers() [2]int64
}

type AnswerMsg struct {
	day    int
	part   int
	answer int64
}

func (d *dayState) calculateCmd(preset int, part int) tea.Cmd {
	out := &d.out[part]
	if out.isCalculating() {
		return nil
	}
	input := d.presets.input(preset)
	outputCh := make(chan []string)
	out.ch = outputCh
	out.answer = nil
	out.lines = nil
	return tea.Batch(
		d.runCalculateCmd(part, input, outputCh),
		d.recvOutputCmd(part, outputCh),
	)
}

func (d dayState) runCalculateCmd(part int, input []string, outputCh chan<- []string) tea.Cmd {
	return func() tea.Msg {
		defer close(outputCh)
		answerMsg := AnswerMsg{
			day:    d.day,
			part:   part,
			answer: 0,
		}
		if part == 0 {
			answerMsg.answer = d.calculate.Part1(input, outputCh)
		} else {
			answerMsg.answer = d.calculate.Part2(input, outputCh)
		}
		return answerMsg
	}
}

type OutputMsg struct {
	day   int
	part  int
	lines []string
}

func (d dayState) recvOutputCmd(part int, outputChan <-chan []string) tea.Cmd {
	return func() tea.Msg {
		lines, ok := <-outputChan
		if !ok {
			return nil
		}
		return OutputMsg{
			day:   d.day,
			part:  part,
			lines: lines,
		}
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
		day06.Calculate{},
		day07.Calculate{},
		day08.Calculate{},
		day09.Calculate{},
		day10.Calculate{},
		day11.Calculate{},
	}
}

func (m *model) scheduleAutosolve() tea.Cmd {
	var commands []tea.Cmd
	for day := range m.dayStates {
		d := &m.dayStates[day]
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
