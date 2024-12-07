package main

import (
	"fmt"
	"image/color"
	"os"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
)

type State int

const (
	CalendarState State = iota
	DayState
	CalculateState
	ExitState
)

type model struct {
	state           State
	size            *tea.WindowSizeMsg
	presets         *loadedPresets
	originalBgColor *color.Color
	selectedDay     int
	inputScroll     [25]int
	selectedPreset  [25]int
	selectedPart    [25]int
	answer          *int64
}

func (m model) Init() (tea.Model, tea.Cmd) {
	for i := range m.selectedPreset {
		m.selectedPreset[i] = 1
	}
	return m, tea.Batch(
		loadPresets,
		tea.Sequence(
			tea.RequestBackgroundColor,
			tea.SetBackgroundColor(lipgloss.Color("#0f0f23")),
		),
	)
}

type ExitMsg struct{}

func exit() tea.Msg {
	return ExitMsg{}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.size = &msg
	case loadedPresets:
		m.presets = &msg
	case tea.BackgroundColorMsg:
		m.originalBgColor = &msg.Color
	case ExitMsg:
		return m, tea.Quit
	case tea.KeyMsg:
		switch msg.String() {
		case "up":
			switch m.state {
			case CalendarState:
				m.selectedDay--
				if m.selectedDay < 0 {
					m.selectedDay = len(m.presets.ascii) - 1
				}
			case DayState:
				if m.inputScroll[m.selectedDay] > 0 {
					m.inputScroll[m.selectedDay]--
				}
			}
		case "down":
			switch m.state {
			case CalendarState:
				m.selectedDay++
				if m.selectedDay >= len(m.presets.ascii) {
					m.selectedDay = 0
				}
			case DayState:
				maxScroll := 0
				if m.size != nil {
					input := m.presets.input(m.selectedDay, m.selectedPreset[m.selectedDay])
					maxScroll = len(input) - m.size.Height + 6
				}
				if m.inputScroll[m.selectedDay] < maxScroll {
					m.inputScroll[m.selectedDay]++
				}
			}
		case "enter":
			if m.state == CalendarState {
				m.state = DayState
				m.answer = nil
			}
		case "left":
			if m.state == DayState {
				m.selectedPart[m.selectedDay] = 0
			}
		case "right":
			if m.state == DayState {
				m.selectedPart[m.selectedDay] = 1
			}
		case "1", "2", "3", "4", "5", "6", "7", "8", "9":
			if m.state == DayState {
				nextNum := int(msg.String()[0] - '0')
				for _, preset := range m.presets.days[m.selectedDay] {
					if preset.num == nextNum {
						m.selectedPreset[m.selectedDay] = nextNum
					}
				}
			}
		case "esc", "q":
			if m.state == CalendarState {
				m.state = ExitState
				return m, exit
			}
			m.state = CalendarState
		case "ctrl+c":
			m.state = ExitState
			return m, exit
		}
	}

	return m, nil
}

var highlightStyle lipgloss.Style = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#00cc00")).
	Padding(1)
var controlStyle lipgloss.Style = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#009900")).
	Padding(1)
var textStyle lipgloss.Style = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#cccccc"))
var calendarSelectedStyle lipgloss.Style = textStyle.
	Background(lipgloss.Color("#24243b"))
var dataStyle lipgloss.Style = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#cccccc")).
	Background(lipgloss.Color("#10101a"))

func (m model) View() string {
	if m.size == nil || m.presets == nil {
		return ""
	}

	if m.state == ExitState && m.originalBgColor != nil {
		// Clear screen before exit
		return lipgloss.NewStyle().
			Height(m.size.Height).
			Width(m.size.Width).
			Background(*m.originalBgColor).
			Render("wgdgqwkduyqgwdku")
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		m.headerView(),
		m.bodyView(),
		m.footerView(),
	)
}

func (m model) headerView() string {
	controls := []string{"title placeholder"}
	if m.state == CalendarState {
		controls = append(controls, controlStyle.Render("[Esc: exit]"))
	} else {
		for _, preset := range m.presets.days[m.selectedDay] {
			controlText := fmt.Sprintf("[%d: %s]", preset.num, preset.tag)
			if m.selectedPreset[m.selectedDay] == preset.num {
				controls = append(controls, highlightStyle.Render(controlText))
			} else {
				controls = append(controls, controlStyle.Render(controlText))
			}
		}
		controls = append(controls, controlStyle.Render("[Esc: back]"))
	}

	titleText := "Advent of Code 2024"
	if m.state != CalendarState {
		titleText = fmt.Sprintf("%s: Day %d", titleText, m.selectedDay+1)
	}
	title := highlightStyle.Render(titleText)

	return m.joinWithGap(
		[]string{title},
		controls,
	)
}

func (m model) bodyView() string {
	if m.state == CalendarState {
		return m.calendarSelectView()
	} else {
		return m.inputAndLogView()
	}
}

func (m model) calendarSelectView() string {
	var calendarLines []string
	for day, asciiLine := range m.presets.ascii {
		asciiLineWithDay := fmt.Sprintf("%s  %2d", asciiLine, day+1)
		if m.selectedDay == day {
			calendarLines = append(calendarLines, calendarSelectedStyle.Render(asciiLineWithDay))
		} else {
			calendarLines = append(calendarLines, textStyle.Render(asciiLineWithDay))
		}
	}
	gapHeight := m.size.Height - len(m.presets.ascii) - 6
	calendarLines = append(calendarLines, textStyle.Height(gapHeight).Render(""))
	return lipgloss.JoinVertical(
		lipgloss.Left,
		calendarLines...,
	)
}

func (m model) inputAndLogView() string {
	input := m.presets.input(m.selectedDay, m.selectedPreset[m.selectedDay])
	scrollTop := m.inputScroll[m.selectedDay]
	scrollBottom := min(scrollTop+m.size.Height-6, len(input)-1)
	window := input[scrollTop:scrollBottom]
	return dataStyle.
		Height(m.size.Height - 6).
		MarginLeft(1).
		Render(strings.Join(window, "\n"))
}

func (m model) footerView() string {
	if m.state == CalendarState {
		return controlStyle.
			Width(m.size.Width - 1).
			Align(lipgloss.Right).
			Render("[Enter: select]")
	}
	answerLabel := textStyle.Padding(1).Render("Answer:")
	answer := dataStyle.Render(m.answerText())
	var controls []string
	for part, partText := range []string{"[←: Part 1]", "[→: Part 2]"} {
		if m.selectedPart[m.selectedDay] == part {
			controls = append(controls, highlightStyle.Render(partText))
		} else {
			controls = append(controls, controlStyle.Render(partText))
		}
	}
	if m.state == CalculateState {
		controls = append(controls, highlightStyle.Render("[ calculating... ]"))
	} else {
		controls = append(controls, controlStyle.Render("[Enter: calculate]"))
	}
	return m.joinWithGap(
		[]string{answerLabel, answer},
		controls,
	)
}

func (m model) answerText() string {
	if m.state == DayState && m.answer != nil {
		return strconv.FormatInt(*m.answer, 10)
	}
	return "-"
}

func (m model) joinWithGap(leftWidgets []string, rightWidgets []string) string {
	widgetsWidth := 0
	for _, widget := range leftWidgets {
		widgetsWidth += lipgloss.Width(widget)
	}
	for _, widget := range rightWidgets {
		widgetsWidth += lipgloss.Width(widget)
	}
	gapWidth := m.size.Width - widgetsWidth - 1
	var widgets []string
	widgets = append(widgets, leftWidgets...)
	widgets = append(widgets, lipgloss.NewStyle().Width(gapWidth).Render(""))
	widgets = append(widgets, rightWidgets...)
	return lipgloss.JoinHorizontal(lipgloss.Center, widgets...)
}

func main() {
	p := tea.NewProgram(model{})
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
