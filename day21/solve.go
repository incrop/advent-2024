package day21

import (
	"incrop/advent-2024/out"
	"iter"
	"log"
	"strconv"
	"strings"
)

type Solve struct{}

func (d Solve) Part1(input []string, outputCh chan<- []string) string {
	complexitySum := int64(0)
	for _, code := range parse(input) {
		c := newChain(
			directionalKeypad(),
			directionalKeypad(),
			numericKeypad(),
		)
		minInput, minLength := c.shortestInputForOutput(code, true)
		if len(minInput) != int(minLength) {
			panic(minLength)
		}
		outputCh <- c.show()
		for _, btn := range minInput {
			c.feed(btn)
			outputCh <- c.show()
		}
		numCodePart, err := strconv.Atoi(code[:len(code)-1])
		if err != nil {
			log.Fatal(err)
		}
		complexitySum += minLength * int64(numCodePart)
	}
	return strconv.FormatInt(complexitySum, 10)
}

func (d Solve) Part2(input []string, outputCh chan<- []string) string {
	keypads := []*keypad{}
	for range 25 {
		keypads = append(keypads, directionalKeypad())
	}
	keypads = append(keypads, numericKeypad())
	l := out.NewLog(outputCh)
	complexitySum := int64(0)
	for _, code := range parse(input) {
		c := newChain(keypads...)
		_, minLength := c.shortestInputForOutput(code, false)
		numCodePart, err := strconv.Atoi(code[:len(code)-1])
		if err != nil {
			log.Fatal(err)
		}
		complexity := int64(numCodePart) * minLength
		l.Printf("%03d * %d = %d", numCodePart, minLength, complexity)
		complexitySum += complexity
	}
	return strconv.FormatInt(complexitySum, 10)
}

type keypad struct {
	i, j    int
	buttons [][]rune
}

const (
	up       = '^'
	right    = '>'
	down     = 'v'
	left     = '<'
	activate = 'A'
)

func numericKeypad() (k *keypad) {
	k = new(keypad)
	k.buttons = [][]rune{
		{'7', '8', '9'},
		{'4', '5', '6'},
		{'1', '2', '3'},
		{' ', '0', 'A'},
	}
	k.i, k.j = 3, 2
	return
}

func directionalKeypad() (k *keypad) {
	k = new(keypad)
	k.buttons = [][]rune{
		{' ', '^', 'A'},
		{'<', 'v', '>'},
	}
	k.i, k.j = 0, 2
	return
}

func (k *keypad) move(d rune) {
	switch d {
	case up:
		k.i--
	case right:
		k.j++
	case down:
		k.i++
	case left:
		k.j--
	}
}

func (k *keypad) activate() rune {
	return k.buttons[k.i][k.j]
}

func (k *keypad) show() (o []string) {
	runes := make([][]rune, len(k.buttons)*2+1)
	for i := range runes {
		runes[i] = make([]rune, len(k.buttons[0])*4+1)
		for j := range runes[i] {
			runes[i][j] = ' '
		}
	}
	for i, row := range k.buttons {
		for j, btn := range row {
			ic, jc := i*2+1, j*4+2
			runes[ic][jc] = btn
			if btn == ' ' {
				continue
			}
			copy(runes[ic-1][jc-2:jc+3], []rune("+---+"))
			runes[ic][jc-2], runes[ic][jc+2] = '|', '|'
			copy(runes[ic+1][jc-2:jc+3], []rune("+---+"))
			if i == k.i && j == k.j {
				runes[ic][jc-1], runes[ic][jc+1] = '[', ']'
			}
		}
	}
	for _, row := range runes {
		o = append(o, string(row))
	}
	return
}

func (k *keypad) inputsBetweenSeq(s, t rune) iter.Seq[string] {
	is, js, it, jt := 0, 0, 0, 0
	for i, row := range k.buttons {
		for j, btn := range row {
			if btn == s {
				is, js = i, j
			}
			if btn == t {
				it, jt = i, j
			}
		}
	}
	return func(yield func(string) bool) {
		interruped := false
		var travel func(path []rune, i, j int)
		travel = func(path []rune, i, j int) {
			if interruped {
				return
			}
			if i == it && j == jt {
				interruped = !yield(string(append(path, activate)))
				return
			}
			if k.buttons[i][j] == ' ' {
				return
			}
			if i != it {
				input, di := 'v', 1
				if i > it {
					input, di = '^', -1
				}
				travel(append(path, input), i+di, j)
			}
			if j != jt {
				input, dj := '>', 1
				if j > jt {
					input, dj = '<', -1
				}
				travel(append(path, input), i, j+dj)
			}
		}
		travel(make([]rune, 0, 6), is, js)
	}
}

type chain struct {
	keypads []*keypad
	input   []rune
	outputs [][]rune
}

func newChain(keypads ...*keypad) (c *chain) {
	c = new(chain)
	c.keypads = keypads
	c.outputs = make([][]rune, len(keypads))
	return
}

func (c *chain) feed(r rune) {
	c.input = append(c.input, r)
	for i, k := range c.keypads {
		switch r {
		case up, right, down, left:
			k.move(r)
			return
		case activate:
			r = k.activate()
			c.outputs[i] = append(c.outputs[i], r)
		default:
			panic(r)
		}
	}
}

func (c *chain) show() (o []string) {
	o = append(o, string(c.input))
	for i, k := range c.keypads {
		o = append(o, k.show()...)
		o = append(o, string(c.outputs[i]))
	}
	return
}

func (c *chain) shortestInputForOutput(output string, collectString bool) (string, int64) {
	type cacheKey struct {
		i          int
		btn1, btn2 rune
	}
	type cacheValue struct {
		input  string
		length int64
	}
	cache := map[cacheKey]cacheValue{}
	var find func(i int, output string, acc int64) cacheValue
	find = func(i int, output string, acc int64) cacheValue {
		if i == -1 {
			return cacheValue{output, acc}
		}
		totalLength := int64(0)
		var sb strings.Builder
		btn1 := 'A'
		for _, btn2 := range output {
			key := cacheKey{i, btn1, btn2}
			min := cacheValue{}
			if cached, ok := cache[key]; ok {
				min = cached
			} else {
				for prevLevelOutput := range c.keypads[i].inputsBetweenSeq(btn1, btn2) {
					value := find(i-1, prevLevelOutput, int64(len(prevLevelOutput)))
					if min.length == 0 || value.length < min.length {
						min = value
					}
				}
				cache[key] = min
			}
			if collectString {
				sb.WriteString(min.input)
			}
			totalLength += min.length
			btn1 = btn2
		}
		return cacheValue{sb.String(), totalLength}
	}
	value := find(len(c.keypads)-1, output, int64(len(output)))
	return value.input, value.length
}

func parse(input []string) (codes []string) {
	return append(codes, input...)
}

func (d Solve) CorrectAnswers() [2]string {
	return [2]string{"123096", "154517692795352"}
}
