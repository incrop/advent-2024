package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func ascii(selectedDay int, dayStars []int) []string {
	return applyStyles(
		selectedDay,
		`                                                 `,
		`          .-----.          .------------------.  `,
		`           66666             $       bbbb        `,
		`       .--'~ ~ ~|        .-' *       \  /     '-.`,
		`        6  g  6             GOG   bbbbbbbbbbb    `,
		`    .--'~  ,* ~ |        |  >o<   \_\_\|_/__/   |`,
		`       6 giiig 6           GRGBG Rbbb   bbbbbbbb `,
		`.---': ~ '(~), ~|        | >@>O< o-_/.()__------|`,
		`                                                 `,
		`|               |        |                      |`,
		`                                                 `,
		`|               |        |          ..          |`,
		`                                                 `,
		`|               |        |        .'  '.        |`,
		`                                                 `,
		`|               |        |        |    |        |`,
	)
}

type styleOverlay struct {
	len   int
	style *lipgloss.Style
}

func applyStyles(selectedDay int, styledAscii ...string) (rendered []string) {
	if len(styledAscii)%2 > 0 {
		panic("Expected even number of rows")
	}
	for i, day := 0, 0; i < len(styledAscii); i, day = i+2, day+1 {
		styleRow := styledAscii[i]
		asciiRow := styledAscii[i+1]
		if len(styleRow) != len(asciiRow) {
			panic("Style and ASCII rows should have same length")
		}
		styleOverlays := parseStyleOverlays(styleRow)
		rendered = append(rendered, applyOverlays(
			styleOverlays,
			asciiRow,
			day,
			0,
			day == selectedDay,
		))
	}
	return
}

var white lipgloss.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#cccccc"))
var beige lipgloss.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#e3b585"))
var brown lipgloss.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#886655"))
var lgreen lipgloss.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#00cc00"))
var green lipgloss.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#009900"))
var indigo lipgloss.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#5555bb"))
var orange lipgloss.Style = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#ff9900"))
var red lipgloss.Style = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#ff0000"))
var blue lipgloss.Style = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#0066ff"))
var star lipgloss.Style = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#ffff66"))

func parseStyleOverlays(styleRow string) (styleOverlays []styleOverlay) {
	currentOverlay := styleOverlay{
		len:   0,
		style: &white,
	}
	for _, r := range styleRow {
		var nextStyle *lipgloss.Style
		switch r {
		case ' ':
			nextStyle = &white
		case '6':
			nextStyle = &beige
		case 'b':
			nextStyle = &brown
		case 'g':
			nextStyle = &lgreen
		case 'G':
			nextStyle = &green
		case 'i':
			nextStyle = &indigo
		case 'O':
			nextStyle = &orange
		case 'R':
			nextStyle = &red
		case 'B':
			nextStyle = &blue
		case '$':
			nextStyle = &star
		default:
			panic(fmt.Sprintf("Unexpected symbol %c", r))
		}
		if nextStyle == currentOverlay.style {
			currentOverlay.len++
		} else {
			styleOverlays = append(styleOverlays, currentOverlay)
			currentOverlay = styleOverlay{
				len:   1,
				style: nextStyle,
			}
		}
	}
	styleOverlays = append(styleOverlays, currentOverlay)
	return
}

func applyOverlays(styleOverlays []styleOverlay, asciiRow string, day int, stars int, selected bool) string {
	var sb strings.Builder
	i := 0
	for _, styleOverlay := range styleOverlays {
		asciiPart := asciiRow[i : i+styleOverlay.len]
		style := *styleOverlay.style
		if selected {
			style = style.Background(lipgloss.Color("#24243b"))
		}
		sb.WriteString(style.Render(asciiPart))
		i += styleOverlay.len
	}
	if day > 0 {
		style := white
		if selected {
			style = style.Background(lipgloss.Color("#24243b"))
		}
		sb.WriteString(style.Render(fmt.Sprintf("  %2d", day)))
		style = star
		if selected {
			style = style.Background(lipgloss.Color("#24243b"))
		}
		sb.WriteString(style.Render(" " + strings.Repeat("*", stars)))
	}
	return sb.String()
}
