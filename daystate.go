package main

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
)

type dayState struct {
	day              int
	inputs           dayInputs
	selectedInput    int
	selectedPart     int
	scrollX, scrollY int
	separatorOffset  int
	solve            Solve
	out              [2]output
}

type output struct {
	ch     <-chan []string
	lines  []string
	answer string
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
	case "-", "_":
		d.separatorOffset--
	case "+", "=":
		d.separatorOffset++
	case "enter", "space":
		return d.solveCmd(d.selectedInput, d.selectedPart)
	case "tab":
		d.selectedPart = 1 - d.selectedPart
	case "1", "2", "3", "4", "5", "6", "7", "8", "9":
		inputNum := int(msg.String()[0] - '0')
		if !d.out[d.selectedPart].isCalculating() && d.inputs.lines(inputNum) != nil {
			d.selectedInput = inputNum
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
	newAnswer := msg.answer
	out := &d.out[msg.part]
	out.answer = newAnswer
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
	for _, input := range d.inputs {
		controlText := fmt.Sprintf("[%d: %s]", input.num, input.tag)
		if d.selectedInput == input.num {
			controls = append(controls, highlightStyle.Render(controlText))
		} else {
			controls = append(controls, controlStyle.Render(controlText))
		}
	}
	for part, partText := range []string{"[%sPart 1]", "[%sPart 2]"} {
		if d.selectedPart == part {
			controls = append(controls, highlightStyle.Render(fmt.Sprintf(partText, "")))
		} else {
			controls = append(controls, controlStyle.Render(fmt.Sprintf(partText, "Tab: ")))
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
	answer := d.answerView()

	var controls []string
	controls = append(controls, controlStyle.Render("[← ↑ → ↓: scroll]"))
	controls = append(controls, controlStyle.Render("[- +: separator]"))
	if d.out[d.selectedPart].isCalculating() {
		controls = append(controls, highlightStyle.Render("[ solving... ]"))
	} else {
		controls = append(controls, controlStyle.Render("[Enter: solve]"))
	}

	return joinHorizontalWithGap(
		[]string{answerLabel, answer},
		controls,
		maxWidth,
	)
}

func (d dayState) answerView() string {
	answer := d.out[d.selectedPart].answer
	if answer == "" {
		return "-"
	}
	view := dataStyle.Render(answer)
	coorectAnswer := d.solve.CorrectAnswers()[d.selectedPart]
	if coorectAnswer == answer {
		view += starStyle.Render("*")
	}
	return view
}

func (d dayState) bodyView(size tea.WindowSizeMsg) string {
	width := (size.Width - 3) / 2
	height := size.Height
	off := d.separatorOffset
	minWidth := 10
	if off < 0 && width+off < minWidth {
		off = minWidth - width
	} else if off > 0 && width-off < minWidth {
		off = width - minWidth
	}

	input := d.inputs.lines(d.selectedInput)
	output := d.out[d.selectedPart].lines
	inputWindow := cropWindow(input, d.scrollX, d.scrollY, width+off, height)
	outputWindow := cropWindow(output, d.scrollX, d.scrollY, width-off, height)

	style := dataStyle.Height(height).MarginLeft(1)
	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		style.Width(width+off).Render(strings.Join(inputWindow, "\n")),
		style.Width(width-off).Render(strings.Join(outputWindow, "\n")),
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
	for _, lineStr := range lines[minY:maxY] {
		line := []rune(lineStr)
		minX := scrollX
		if minX > len(line) {
			minX = len(line)
		}
		maxX := minX + width
		if maxX > len(line) {
			maxX = len(line)
		}
		cropLine := line[minX:maxX]
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
