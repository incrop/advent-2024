package dayXX

type Solve struct{}

func (d Solve) Part1(input []string, outputCh chan<- []string) string {
	p := parse(input)
	outputCh <- p.output()
	return ""
}

func (d Solve) Part2(input []string, outputCh chan<- []string) string {
	p := parse(input)
	outputCh <- p.output()
	return ""
}

type parsed []string

func parse(input []string) (parsed parsed) {
	for _, line := range input {
		parsed = append(parsed, line)
	}
	return
}

func (p parsed) output() (lines []string) {
	return p
}

func (d Solve) CorrectAnswers() [2]string {
	return [2]string{"not solved", "not solved"}
}
