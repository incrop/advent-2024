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
	"incrop/advent-2024/day12"
	"incrop/advent-2024/day13"
	"incrop/advent-2024/day14"
	"incrop/advent-2024/day15"
	"incrop/advent-2024/day16"
	"incrop/advent-2024/day17"
	"incrop/advent-2024/day18"

	tea "github.com/charmbracelet/bubbletea/v2"
)

type Solve interface {
	Part1(input []string, outputCh chan<- []string) int64
	Part2(input []string, outputCh chan<- []string) int64
	CorrectAnswers() [2]int64
}

type AnswerMsg struct {
	day    int
	part   int
	answer int64
}

func (d *dayState) solveCmd(inputNum int, part int) tea.Cmd {
	out := &d.out[part]
	if out.isCalculating() {
		return nil
	}
	input := d.inputs.lines(inputNum)
	outputCh := make(chan []string)
	out.ch = outputCh
	out.answer = nil
	out.lines = nil
	return tea.Batch(
		d.runSolveCmd(part, input, outputCh),
		d.recvOutputCmd(part, outputCh),
	)
}

func (d dayState) runSolveCmd(part int, input []string, outputCh chan<- []string) tea.Cmd {
	return func() tea.Msg {
		defer close(outputCh)
		answerMsg := AnswerMsg{
			day:    d.day,
			part:   part,
			answer: 0,
		}
		if part == 0 {
			answerMsg.answer = d.solve.Part1(input, outputCh)
		} else {
			answerMsg.answer = d.solve.Part2(input, outputCh)
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

func collectCalculations() [26]Solve {
	return [26]Solve{
		nil,
		day01.Solve{},
		day02.Solve{},
		day03.Solve{},
		day04.Solve{},
		day05.Solve{},
		day06.Solve{},
		day07.Solve{},
		day08.Solve{},
		day09.Solve{},
		day10.Solve{},
		day11.Solve{},
		day12.Solve{},
		day13.Solve{},
		day14.Solve{},
		day15.Solve{},
		day16.Solve{},
		day17.Solve{},
		day18.Solve{},
	}
}

func (m *model) scheduleAutosolve() tea.Cmd {
	var commands []tea.Cmd
	for day := range m.dayStates {
		d := &m.dayStates[day]
		if d.solve == nil {
			continue
		}
		commands = append(
			commands,
			d.solveCmd(1, 0),
			d.solveCmd(1, 1),
		)
	}
	return tea.Batch(commands...)
}
