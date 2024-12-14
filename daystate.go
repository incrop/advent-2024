package main

import (
	"fmt"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
)

type dayState struct {
	day              int
	presets          dayPresets
	selectedPreset   int
	selectedPart     int
	scrollX, scrollY int
	calculate        Calculate
	out              [2]output
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
		if d.scrollY > 0 {
			d.scrollY--
		}
	case "down":
		d.scrollY++
	case "left":
		if d.scrollX > 0 {
			d.scrollX--
		}
	case "right":
		d.scrollX++
	case "enter", "space":
		return d.calculateCmd(d.selectedPreset, d.selectedPart)
	case "tab":
		d.selectedPart = 1 - d.selectedPart
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
	for part, partText := range []string{"[%sPart 1]", "[%sPart 2]"} {
		if d.selectedPart == part {
			controls = append(controls, highlightStyle.Render(fmt.Sprintf(partText, "")))
		} else {
			controls = append(controls, controlStyle.Render(fmt.Sprintf(partText, "Tab: ")))
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
	width := (size.Width - 3) / 2
	height := size.Height

	input := d.presets.input(d.selectedPreset)
	output := d.out[d.selectedPart].lines
	inputWindow := cropWindow(input, d.scrollX, d.scrollY, width, height)
	outputWindow := cropWindow(output, d.scrollX, d.scrollY, width, height)

	style := dataStyle.Width(width).Height(height).MarginLeft(1)
	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		style.Render(strings.Join(inputWindow, "\n")),
		style.Render(strings.Join(outputWindow, "\n")),
	)
}

func cropWindow(lines []string, scrollX, scrollY, width, height int) (window []string) {
	minY := scrollY
	if minY > len(lines) {
		minY = len(lines)
	}
	maxY := minY + height
	if maxY > len(lines) {
		maxY = len(lines)
	}
	for _, line := range lines[minY:maxY] {
		minX := scrollX
		if minX > len(line) {
			minX = len(line)
		}
		maxX := minX + width
		if maxX > len(line) {
			maxX = len(line)
		}
		cropLine := []rune(line[minX:maxX])
		if minX > 0 {
			if len(cropLine) > 0 {
				cropLine[0] = '←'
			} else {
				cropLine = append(cropLine, '←')
			}
		}
		if maxX < len(line) && len(cropLine) > 0 {
			cropLine[len(cropLine)-1] = '→'
		}
		window = append(window, string(cropLine))
	}
	if minY > 0 {
		if len(window) > 0 {
			window[0] = " ↑↑↑"
		} else {
			window = append(window, "↑↑↑")
		}
	}
	if maxY < len(lines) && len(window) > 0 {
		window[len(window)-1] = " ↓↓↓"
	}
	return
}
