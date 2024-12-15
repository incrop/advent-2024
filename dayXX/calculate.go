package dayXX

type Calculate struct{}

func (d Calculate) Part1(input []string, outputCh chan<- []string) int64 {
	p := parse(input)
	outputCh <- p.output()
	return 0
}

func (d Calculate) Part2(input []string, outputCh chan<- []string) int64 {
	p := parse(input)
	outputCh <- p.output()
	return 0
}

type parsed []string

func parse(input []string) (parsed parsed) {
	for _, line := range input {
		parsed = append(parsed, line)
	}
	return
}

func (p parsed) output() (lines []string) {
	return
}

func (d Calculate) CorrectAnswers() [2]int64 {
	return [2]int64{-1, -1}
}
