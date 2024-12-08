package main

import (
	"fmt"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
)

type dayState struct {
	day            int
	presets        dayPresets
	selectedPreset int
	selectedPart   int
	inputScroll    int
	calculate      Calculate
	out            [2]output
}

type output struct {
	ch     <-chan []string
	lines  []string
	answer *int64
}

func (out output) isCalculating() bool {
	return out.ch != nil
}

func (d *dayState) handleKeyMsg(msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {
	case "up":
		if d.inputScroll > 0 {
			d.inputScroll--
		}
	case "down":
		d.inputScroll++
	case "enter", "space":
		return d.calculateCmd(d.selectedPreset, d.selectedPart)
	case "left":
		d.selectedPart = 0
	case "right":
		d.selectedPart = 1
	case "1", "2", "3", "4", "5", "6", "7", "8", "9":
		presetNum := int(msg.String()[0] - '0')
		if !d.out[d.selectedPart].isCalculating() && d.presets.input(presetNum) != nil {
			d.selectedPreset = presetNum
		}
	}
	return nil
}

func (d *dayState) handleOutputMsg(msg OutputMsg) tea.Cmd {
	out := &d.out[msg.part]
	out.lines = msg.lines
	return d.recvOutputCmd(msg.part, out.ch)
}

func (d *dayState) handleAnswerMsg(msg AnswerMsg) tea.Cmd {
	newAnswer := int64(msg.answer)
	out := &d.out[msg.part]
	out.answer = &newAnswer
	out.ch = nil
	return nil
}

func (d dayState) view(size tea.WindowSizeMsg) string {
	header := d.headerView(size.Width)
	footer := d.footerView(size.Width)
	bodySize := size
	bodySize.Height -= lipgloss.Height(header) + lipgloss.Height(footer)
	body := d.bodyView(bodySize)
	return lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		body,
		footer,
	)
}

func (d dayState) headerView(maxWidth int) string {
	titleText := fmt.Sprintf("Advent of Code 2024: Day %d", d.day)
	title := highlightStyle.Render(titleText)

	var controls []string
	for _, preset := range d.presets {
		controlText := fmt.Sprintf("[%d: %s]", preset.num, preset.tag)
		if d.selectedPreset == preset.num {
			controls = append(controls, highlightStyle.Render(controlText))
		} else {
			controls = append(controls, controlStyle.Render(controlText))
		}
	}
	controls = append(controls, controlStyle.Render("[Esc: back]"))

	return joinHorizontalWithGap(
		[]string{title},
		controls,
		maxWidth,
	)
}

func (d dayState) footerView(maxWidth int) string {
	answerLabel := textStyle.Padding(1).Render("Answer:")
	answer := dataStyle.Render(d.answerText())

	var controls []string
	for part, partText := range []string{"[←: Part 1]", "[→: Part 2]"} {
		if d.selectedPart == part {
			controls = append(controls, highlightStyle.Render(partText))
		} else {
			controls = append(controls, controlStyle.Render(partText))
		}
	}
	if d.out[d.selectedPart].isCalculating() {
		controls = append(controls, highlightStyle.Render("[ calculating... ]"))
	} else {
		controls = append(controls, controlStyle.Render("[Enter: calculate]"))
	}

	return joinHorizontalWithGap(
		[]string{answerLabel, answer},
		controls,
		maxWidth,
	)
}

func (d dayState) answerText() string {
	answer := d.out[d.selectedPart].answer
	if answer == nil {
		return "-"
	}
	return strconv.FormatInt(*answer, 10)
}

func (d dayState) bodyView(size tea.WindowSizeMsg) string {
	input := d.presets.input(d.selectedPreset)
	scrollBottom := min(d.inputScroll+size.Height, len(input))
	window := input[d.inputScroll:scrollBottom]
	style := dataStyle.
		Width((size.Width - 1) / 2).
		Height(size.Height).
		MarginLeft(1)

	outLines := d.out[d.selectedPart].lines
	if len(outLines) > size.Height {
		outLines = outLines[:size.Height]
	}
	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		style.Render(strings.Join(window, "\n")),
		style.Render(strings.Join(outLines, "\n")),
	)
}
