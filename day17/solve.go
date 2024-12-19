package day17

import (
	"fmt"
	"log"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

type Solve struct{}

func (d Solve) Part1(input []string, outputCh chan<- []string) int64 {
	c := parse(input)
	outputCh <- c.show()
	res := c.runToHalt()
	outputCh <- c.show()
	return res
}

func (d Solve) Part2(input []string, outputCh chan<- []string) int64 {
	c0 := parse(input)
	var find func(currA int64) int64
	find = func(currA int64) int64 {
		for digit := range 10 {
			nextA := currA<<3 + int64(digit)
			c1 := c0.clone()
			c1.a = nextA
			c1.runToHalt()
			if !slices.Equal(c1.output, c1.program[len(c0.program)-len(c1.output):]) {
				continue
			}
			if len(c1.program) == len(c1.output) {
				return nextA
			}
			if solution := find(nextA); solution != 0 {
				return solution
			}
		}
		return 0
	}
	res := find(0)
	if res > 0 {
		c0.a = res
		outputCh <- c0.show()
	}
	return res
}

func (c *computer) runToHalt() (finalOutput int64) {
	for {
		res, tgt, halt := c.calcInstruction()
		if halt {
			break
		}
		if tgt != nil {
			*tgt = res
		}
		if opcode(c.program[c.pointer]) == out {
			c.output = append(c.output, byte(res))
		}
		if opcode(c.program[c.pointer]) == jnz && res != -1 {
			c.pointer = int(res)
			continue
		}
		c.pointer += 2
	}
	for _, b := range c.output {
		finalOutput *= 10
		finalOutput += int64(b)
	}
	return
}

func (c1 *computer) clone() (c2 *computer) {
	c2 = new(computer)
	*c2 = *c1
	c2.program = append([]byte{}, c1.program...)
	c2.output = append([]byte{}, c1.output...)
	return
}

type computer struct {
	a, b, c int64
	pointer int
	program []byte
	output  []byte
}

type opcode byte

const (
	adv opcode = iota
	bxl
	bst
	jnz
	bxc
	out
	bdv
	cdv
)

var registerRegexp = regexp.MustCompile(`^Register ([A|B|C]): (-?\d+)$`)
var programRegexp = regexp.MustCompile(`^Program: ((:?\d,)*\d)$`)

func parse(input []string) (c *computer) {
	c = new(computer)
	for _, line := range input {
		if line == "" {
			continue
		}
		if match := registerRegexp.FindStringSubmatch(line); match != nil {
			val, err := strconv.ParseInt(match[2], 10, 64)
			if err != nil {
				log.Fatal(err)
			}
			switch match[1] {
			case "A":
				c.a = val
			case "B":
				c.b = val
			case "C":
				c.c = val
			}
			continue
		}
		if match := programRegexp.FindStringSubmatch(line); match != nil {
			for _, r := range match[1] {
				if r >= '0' && r <= '9' {
					c.program = append(c.program, byte(r-'0'))
				}
			}
		}
	}
	return
}

func (c *computer) show() (lines []string) {
	lines = append(lines, fmt.Sprintf("Register A: 0o%s", strconv.FormatInt(c.a, 8)))
	lines = append(lines, fmt.Sprintf("Register B: 0o%s", strconv.FormatInt(c.b, 8)))
	lines = append(lines, fmt.Sprintf("Register C: 0o%s", strconv.FormatInt(c.c, 8)))
	lines = append(lines, strings.Repeat(" ", c.pointer*2)+"↓")
	programStr := make([]byte, len(c.program)*2)
	for i, b := range c.program {
		programStr[i*2] = '0' + b
		programStr[i*2+1] = ','
	}
	lines = append(lines, string(programStr[:len(programStr)-1]))
	lines = append(lines, "")
	if len(c.output) == 0 {
		lines = append(lines, "out: ")
	} else {
		outputStr := make([]byte, len(c.output)*2)
		for i, b := range c.output {
			outputStr[i*2] = '0' + b
			outputStr[i*2+1] = ','
		}
		lines = append(lines, "out: "+string(outputStr[:len(outputStr)-1]))
	}
	lines = append(lines, "")
	lines = append(lines, c.showInstruction())
	return lines
}

func (c *computer) combo(b byte) int64 {
	switch b {
	case 0, 1, 2, 3:
		return int64(b)
	case 4:
		return c.a
	case 5:
		return c.b
	case 6:
		return c.c
	default:
		panic(b)
	}
}

func showCombo(b byte) rune {
	switch b {
	case 0, 1, 2, 3:
		return rune('0' + b)
	case 4:
		return 'A'
	case 5:
		return 'B'
	case 6:
		return 'C'
	default:
		panic(b)
	}
}

func (c *computer) showInstruction() string {
	res, _, halt := c.calcInstruction()
	if halt {
		return "halt"
	}
	inst := opcode(c.program[c.pointer])
	oper := c.program[c.pointer+1]
	switch inst {
	case adv:
		return fmt.Sprintf("adv: A / 1<<%c = 0o%s -> A", showCombo(oper), strconv.FormatInt(res, 8))
	case bxl:
		return fmt.Sprintf("bxl: B ^ %d = %d -> B", oper, res)
	case bst:
		return fmt.Sprintf("bst: %c %% 8 = %d -> B", showCombo(oper), res)
	case jnz:
		return fmt.Sprintf("jnz: A -?-> %d", res)
	case bxc:
		return fmt.Sprintf("bxc: B ^ C = %d -> B", res)
	case out:
		return fmt.Sprintf("out: %c %% 8 = 0o%s -> out", showCombo(oper), strconv.FormatInt(res, 8))
	case bdv:
		return fmt.Sprintf("bdv: A / 1<<%c = 0o%s -> B", showCombo(oper), strconv.FormatInt(res, 8))
	case cdv:
		return fmt.Sprintf("cdv: A / 1<<%c = 0o%s -> C", showCombo(oper), strconv.FormatInt(res, 8))
	}
	panic("idk")
}

func (c *computer) calcInstruction() (res int64, tgt *int64, halt bool) {
	if c.pointer >= len(c.program)-1 {
		return 0, nil, true
	}
	inst := opcode(c.program[c.pointer])
	oper := c.program[c.pointer+1]
	switch inst {
	case adv:
		return c.a / (1 << c.combo(oper)), &c.a, false
	case bxl:
		return c.b ^ int64(oper), &c.b, false
	case bst:
		return c.combo(oper) % 8, &c.b, false
	case jnz:
		if c.a == 0 {
			return -1, nil, false
		}
		return int64(oper), nil, false
	case bxc:
		return c.b ^ c.c, &c.b, false
	case out:
		return c.combo(oper) % 8, nil, false
	case bdv:
		return c.a / (1 << c.combo(oper)), &c.b, false
	case cdv:
		return c.a / (1 << c.combo(oper)), &c.c, false
	}
	panic("idk")
}

func (d Solve) CorrectAnswers() [2]int64 {
	return [2]int64{236216121, 90938893795561}
}
