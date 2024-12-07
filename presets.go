package main

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea/v2"
)

type preset struct {
	num   int
	tag   string
	lines []string
}

type dayPresets []preset

type loadedPresets struct {
	ascii []string
	days  [25]dayPresets
}

var presetFileRegexp = regexp.MustCompile(`^day(\d\d)-(\d)-([a-z]+).txt$`)

func loadPresets() tea.Msg {
	var lp loadedPresets
	lp.ascii = loadLines("ascii.txt")
	entries, err := os.ReadDir("./presets")
	if err != nil {
		log.Fatal(err)
	}
	for _, e := range entries {
		match := presetFileRegexp.FindStringSubmatch(e.Name())
		if match == nil {
			continue
		}
		day, err := strconv.Atoi(match[1])
		if err != nil {
			log.Fatal(err)
		}
		num, err := strconv.Atoi(match[2])
		if err != nil {
			log.Fatal(err)
		}
		tag := match[3]
		lp.days[day] = append(lp.days[day], preset{
			num,
			tag,
			loadLines(e.Name()),
		})
	}
	return lp
}

func loadLines(fileName string) (lines []string) {
	file, err := os.Open("./presets/" + fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, strings.TrimSpace(scanner.Text()))
	}
	for len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[0 : len(lines)-1]
	}
	return
}
