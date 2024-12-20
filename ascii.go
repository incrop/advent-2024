package main

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/charmbracelet/lipgloss/v2"
)

func ascii(selectedDay int, dayStars [26]int) []string {
	return applyStyles(
		selectedDay,
		dayStars,
		`                                                 `,
		`          .-----.          .------------------.  `,
		`           66666             $       bbbb        `,
		`       .--'~ ~ ~|        .-' *       \  /     '-.`,
		`        6  g  6             GOG   bbbbbbbbbbb    `,
		`    .--'~  ,* ~ |        |  >o<   \_\_\|_/__/   |`,
		`     6 6 giiig 6           GRGBG Rbbb   bbbbbbbb `,
		`.---': ~ '(~), ~|        | >@>O< o-_/.()__------|`,
		` 9cc96666 ggg 6           GBGOGRG bbbbb       gg `,
		`|@..#'. ~ " ' ~ |        |>O>o<@< \____       .'|`,
		` gccc996666666 $                        bb ggg$  `,
		`|_.~._@'.. ~ ~ *|        | _| |_    ..\_\_ ..'* |`,
		`      9g 9666666          bbb           ggggg gg `,
		`| ||| #@ @'''...|        |...     .'  '.'''../..|`,
		` 9   9g9gg9999               bbb                 `,
		`|#~~~#@@#@ @ #  |        |/\ ''.  |    |   -/  :|`,
		` ccc     rrrrr            $  R  bb               `,
		`|~~..--. _____  |        |* /~\ '.|    | - /  .'|`,
		`        rrrrrrrr          R  $   b        $      `,
		`'---'  ||[][]_\-|        |~/ * \ :|    |  *..'  |`,
		`        rrrrrrr$                bb         ccccc `,
		`       |-------*|        |   /\ .'|    |'''~~~~~|`,
		`        lllllllR             bbb             c   `,
		`       |.......||        |/\ ..'  |    |   . .  |`,
		`                          bbbrr$rr           c c `,
		`       |  -  -  |        |''':::::|    |  .  . .|`,
		`                          r  rr&rr             c `,
		`       |'. -   -|        |.  :::::|    |  .'   .|`,
		`        ggg               ¢  r&&&r      OR  RObb `,
		`       |...'..''|        |☄  :::::|    |..|\..''|`,
		`        g  ggg            r  &&&&&      rrrrrrrr `,
		`       |.  ''.  |        |.  :::::|    |──┬┴┴┴┬─|`,
		`         ggi  gg           r ¢rbrr      rrrrrrrr `,
		`       | '.~  '.|        | : .::::|    |──┤AoC├o|`,
		`        g $ggi g            rr          rrrrrrrr `,
		`       |. *'.~ :|        |  '.    |    |┬o┤ten├─|`,
		`         ggggggg             rB         r$rrrrrr `,
		`       | '..' .'|        |   'o   |    |┘*┤yrs├─|`,
		`         i ggg            r   rr        rrrrrrrr `,
		`       | ~ ..'  |        |:   '.  |    |─┘├┬┬┬┴─|`,
		`        ggg                ¢  $ r       rrrGGOGG `,
		`       |'''))   |        | o. * : '.  .'──┘>>o<<|`,
		`                                                 `,
		`.------'        '------. |          ''          |`,
	)
}

type styleOverlay struct {
	len   int
	style *lipgloss.Style
}

func applyStyles(selectedDay int, dayStars [26]int, styledAscii ...string) (rendered []string) {
	if len(styledAscii)%2 > 0 {
		panic("Expected even number of rows")
	}
	for i, day := 0, 0; i < len(styledAscii); i, day = i+2, day+1 {
		styleRow := styledAscii[i]
		asciiRow := styledAscii[i+1]
		if utf8.RuneCountInString(styleRow) != utf8.RuneCountInString(asciiRow) {
			panic("Style and ASCII rows should have same length")
		}
		styleOverlays := parseStyleOverlays(styleRow)
		rendered = append(rendered, applyOverlays(
			styleOverlays,
			asciiRow,
			day,
			dayStars[day],
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
var dgreen lipgloss.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#01461f"))
var indigo lipgloss.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#5555bb"))
var cyan lipgloss.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#5eabb4"))
var gray lipgloss.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#666666"))
var lavender lipgloss.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#ccccff"))
var orange lipgloss.Style = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#ff9900"))
var red lipgloss.Style = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#ff0000"))
var blue lipgloss.Style = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#0066ff"))
var greener lipgloss.Style = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#009900"))
var star lipgloss.Style = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#ffff66"))
var comet lipgloss.Style = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#bee4e7"))

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
		case '9':
			nextStyle = &dgreen
		case 'i':
			nextStyle = &indigo
		case 'c':
			nextStyle = &cyan
		case 'r':
			nextStyle = &gray
		case 'l':
			nextStyle = &lavender
		case 'O':
			nextStyle = &orange
		case 'R':
			nextStyle = &red
		case 'B':
			nextStyle = &blue
		case '&':
			nextStyle = &greener
		case '$':
			nextStyle = &star
		case '¢':
			nextStyle = &comet
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

func applyOverlays(styleOverlays []styleOverlay, ascii string, day int, stars int, selected bool) string {
	var sb strings.Builder
	sb.WriteString(white.Render(" "))
	i := 0
	asciiRow := []rune(ascii)
	for _, styleOverlay := range styleOverlays {
		asciiPart := asciiRow[i : i+styleOverlay.len]
		style := *styleOverlay.style
		if selected {
			style = style.Background(lipgloss.Color("#24243b"))
		}
		sb.WriteString(style.Render(string(asciiPart)))
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
		sb.WriteString(style.Render(" " + strings.Repeat("*", stars) + strings.Repeat(" ", 2-stars)))
	}
	return sb.String()
}
