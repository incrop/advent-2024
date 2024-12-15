package main

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"strconv"

	tea "github.com/charmbracelet/bubbletea/v2"
)

type dayInput struct {
	num   int
	tag   string
	lines []string
}

type dayInputs []dayInput

func (dis dayInputs) lines(num int) []string {
	for _, input := range dis {
		if input.num == num {
			return input.lines
		}
	}
	return nil
}

type loadedInputs [26]dayInputs

var inputFileRegexp = regexp.MustCompile(`^day(\d\d)-(\d)-([a-z]+).txt$`)

func loadInputs() tea.Msg {
	var inputs loadedInputs
	entries, err := os.ReadDir("./inputs")
	if err != nil {
		log.Fatal(err)
	}
	for _, e := range entries {
		match := inputFileRegexp.FindStringSubmatch(e.Name())
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
		inputs[day] = append(inputs[day], dayInput{
			num,
			tag,
			loadLines(e.Name()),
		})
	}
	return inputs
}

func loadLines(fileName string) (lines []string) {
	file, err := os.Open("./inputs/" + fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	for len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[0 : len(lines)-1]
	}
	return
}
