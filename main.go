package main

import (
	"flag"
	"fmt"
	"image/color"
	"os"

	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
)

type State int

const (
	CalendarState State = iota
	DayState
	ExitState
)

type model struct {
	size            *tea.WindowSizeMsg
	originalBgColor *color.Color
	state           State
	selectedDay     int
	dayStates       [26]dayState
	autosolve       bool
}

func (m model) Init() (tea.Model, tea.Cmd) {
	return m, tea.Batch(
		loadInputs,
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
	case loadedInputs:
		for i, inputs := range &msg {
			m.dayStates[i].inputs = inputs
		}
		if m.autosolve {
			cmd := m.scheduleAutosolve()
			return m, cmd
		}
	case tea.BackgroundColorMsg:
		m.originalBgColor = &msg.Color
	case OutputMsg:
		cmd := m.dayStates[msg.day].handleOutputMsg(msg)
		return m, cmd
	case AnswerMsg:
		cmd := m.dayStates[msg.day].handleAnswerMsg(msg)
		return m, cmd
	case ExitMsg:
		return m, tea.Quit
	case tea.KeyMsg:

		switch msg.String() {
		case "ctrl+c":
			m.state = ExitState
			return m, exit
		case "esc", "q":
			switch m.state {
			case CalendarState:
				m.state = ExitState
				return m, exit
			case DayState:
				m.state = CalendarState
			}
		}
		if m.state == DayState {
			dayCmd := m.dayStates[m.selectedDay].handleKeyMsg(msg)
			return m, dayCmd
		}
		switch msg.String() {
		case "up":
			m.selectedDay--
			if m.selectedDay <= 0 {
				m.selectedDay = 1
				for m.dayStates[m.selectedDay+1].solve != nil {
					m.selectedDay++
				}
			}
		case "down":
			m.selectedDay++
			if m.dayStates[m.selectedDay].solve == nil {
				m.selectedDay = 1
			}
		case "enter", "space":
			m.state = DayState
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
var dataStyle lipgloss.Style = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#cccccc")).
	Background(lipgloss.Color("#10101a"))
var starStyle lipgloss.Style = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#ffff66")).
	MarginLeft(1)

func (m model) View() string {
	if m.size == nil {
		return ""
	}

	if m.state == ExitState && m.originalBgColor != nil {
		// Clear screen before exit
		return lipgloss.NewStyle().
			Height(m.size.Height).
			Width(m.size.Width).
			Background(*m.originalBgColor).
			Render("")
	}

	if m.state == DayState {
		return m.dayStates[m.selectedDay].view(*m.size)
	}

	header := m.headerView()
	footer := m.footerView()
	bodySize := *m.size
	bodySize.Height -= lipgloss.Height(header) + lipgloss.Height(footer)
	body := m.bodyView(bodySize)
	return lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		body,
		footer,
	)
}

func (m model) headerView() string {
	title := highlightStyle.Render("Advent of Code 2024")
	exit := controlStyle.Render("[Esc: exit]")

	return joinHorizontalWithGap(
		[]string{title},
		[]string{exit},
		m.size.Width,
	)
}

func (m model) footerView() string {
	return controlStyle.
		Width(m.size.Width - 1).
		Align(lipgloss.Right).
		Render("[Enter: select]")
}

func (m model) bodyView(size tea.WindowSizeMsg) string {
	calendarLines := ascii(m.selectedDay, m.countDayStars())
	gapHeight := size.Height - len(calendarLines)
	calendarLines = append(calendarLines, textStyle.Height(gapHeight).Render(""))
	return lipgloss.JoinVertical(
		lipgloss.Left,
		calendarLines...,
	)
}

func (m model) countDayStars() (stars [26]int) {
	for day, d := range m.dayStates {
		if d.solve == nil {
			continue
		}
		for part, correctAnswer := range d.solve.CorrectAnswers() {
			answer := d.out[part].answer
			if answer != nil && *answer == correctAnswer {
				stars[day]++
			}
		}
	}
	return
}

func joinHorizontalWithGap(leftWidgets []string, rightWidgets []string, maxWidth int) string {
	widgetsWidth := 0
	for _, widget := range leftWidgets {
		widgetsWidth += lipgloss.Width(widget)
	}
	for _, widget := range rightWidgets {
		widgetsWidth += lipgloss.Width(widget)
	}
	gapWidth := maxWidth - widgetsWidth - 1
	var widgets []string
	widgets = append(widgets, leftWidgets...)
	if gapWidth > 0 {
		widgets = append(widgets, lipgloss.NewStyle().Width(gapWidth).Render(""))
	}
	widgets = append(widgets, rightWidgets...)
	return lipgloss.JoinHorizontal(lipgloss.Center, widgets...)
}

func main() {
	day := flag.Int("day", 0, "from 1 to 25")
	part := flag.Int("part", 1, "1 or 2")
	input := flag.Int("input", 1, "number of input")
	autosolve := flag.Bool("autosolve", true, "run calculations at startup")
	flag.Parse()

	state := CalendarState
	if *day > 0 {
		state = DayState
	} else {
		*day = 1
	}
	p := tea.NewProgram(initModel(state, *day, *part-1, *input, *autosolve))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

func initModel(state State, selectedDay, part, inputNum int, autosolve bool) (m model) {
	m.state = state
	m.selectedDay = selectedDay
	calculations := collectCalculations()
	for day := range m.dayStates {
		d := &m.dayStates[day]
		d.day = day
		d.selectedInput = 1
		d.solve = calculations[day]
	}
	selectedDayState := &m.dayStates[selectedDay]
	selectedDayState.selectedInput = inputNum
	selectedDayState.selectedPart = part
	if state == CalendarState {
		m.autosolve = autosolve
	}
	return
}
