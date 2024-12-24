package day24

import (
	"fmt"
	"incrop/advent-2024/out"
	"regexp"
	"slices"
	"sort"
	"strconv"
	"strings"
)

type Solve struct{}

func (d Solve) Part1(input []string, outputCh chan<- []string) string {
	b := parse(input)
	outNumber := b.run()
	outputCh <- b.output()
	return strconv.FormatInt(outNumber, 10)
}

func (d Solve) Part2(input []string, outputCh chan<- []string) string {
	b := parse(input)
	// Cheated a bit:
	// 1. Find digits where addition breaks:
	//   x = 0000000000000000000000000000000000000000000000 (0)
	//   y = 0111111111111111111111111111111111111111111111 (35184372088831)
	//   z = 0111111000011111111111111111100000001111111111 (34668975883263)
	// 2. Dump the state into graphviz dot format
	// 3. Convert to png and manually check the issue
	swaps := []string{"z10", "vcf", "z17", "fhg", "fsq", "dvb", "z39", "tnc"}
	for i := 0; i < len(swaps); i += 2 {
		b.swap(swaps[i], swaps[i+1])
	}
	l := out.NewLog(outputCh)
	b.calculate(0, 1<<45-1, l)
	l.Append(b.outputGraphviz())
	sort.Strings(swaps)
	return strings.Join(swaps, ",")
}

type operator byte

const (
	and operator = iota + 1
	or
	xor
)

type gate struct {
	in1, in2, out string
	op            operator
	val           int8
}

type board struct {
	inputs, outputs map[string]int8
	gates           []gate
}

func (b board) calculate(x, y int64, l *out.Log) {
	b.setInputs("x", x)
	b.setInputs("y", y)
	z := b.run()
	l.Printf("x = %046b (%d)", x, x)
	l.Printf("y = %046b (%d)", y, y)
	l.Printf("z = %046b (%d)", z, z)
	l.Printf("")
}

func (b board) setInputs(prefix string, value int64) int {
	for i := range 64 {
		inName := fmt.Sprintf("%s%02d", prefix, i)
		inBit := int8(value & 1)
		if _, ok := b.inputs[inName]; !ok {
			return i
		}
		b.inputs[inName] = inBit
		value = value >> 1
	}
	return 64
}

func (b board) run() (outNumber int64) {
	vals := map[string]int8{}
	for in, val := range b.inputs {
		vals[in] = val
	}
	for i, g := range b.gates {
		in1, ok := vals[g.in1]
		if !ok {
			panic(g)
		}
		in2, ok := vals[g.in2]
		if !ok {
			panic(g)
		}
		var out int8
		switch g.op {
		case and:
			out = in1 & in2
		case or:
			out = in1 | in2
		case xor:
			out = in1 ^ in2
		}
		b.gates[i].val = out
		vals[g.out] = out
	}
	for out := range b.outputs {
		b.outputs[out] = vals[out]
	}
	for i := range 64 {
		outName := fmt.Sprintf("z%02d", i)
		outBit, ok := b.outputs[outName]
		if !ok {
			break
		}
		if outBit == 1 {
			outNumber |= 1 << i
		}
	}
	return
}

func (b *board) swap(out1, out2 string) {
	i1, i2 := -1, -1
	for i, g := range b.gates {
		if g.out == out1 {
			i1 = i
		}
		if g.out == out2 {
			i2 = i
		}
	}
	if i1 == -1 || i2 == -1 {
		return
	}
	b.gates[i1].out, b.gates[i2].out = b.gates[i2].out, b.gates[i1].out
	sortGates(b.gates)
}

func parse(input []string) (b board) {
	inputs, nextIdx := parseInputs(input)
	gates, outputs := parseGates(input[nextIdx:])
	return board{
		inputs:  inputs,
		outputs: outputs,
		gates:   gates,
	}
}

var inputRegexp = regexp.MustCompile(`^(.+): (0|1)$`)

func parseInputs(input []string) (inputs map[string]int8, nextIdx int) {
	inputs = map[string]int8{}
	for i, line := range input {
		if line == "" {
			return inputs, i + 1
		}
		match := inputRegexp.FindStringSubmatch(line)
		if match == nil {
			panic(line)
		}
		inputs[match[1]] = int8(match[2][0] - '0')
	}
	return inputs, len(input)
}

var gateRegexp = regexp.MustCompile(`^(.+) (AND|OR|XOR) (.+) -> (.+)$`)

func parseGates(input []string) (gates []gate, outputs map[string]int8) {
	outputs = map[string]int8{}
	for _, line := range input {
		match := gateRegexp.FindStringSubmatch(line)
		if match == nil {
			panic(line)
		}
		var op operator
		switch match[2] {
		case "AND":
			op = and
		case "OR":
			op = or
		case "XOR":
			op = xor
		}
		g := gate{
			in1: match[1],
			in2: match[3],
			out: match[4],
			op:  op,
			val: 0,
		}
		gates = append(gates, g)
		if g.out[0] == 'z' {
			outputs[g.out] = 0
		}
	}
	sortGates(gates)
	return
}

func sortGates(gates []gate) {
	gateByOut := map[string]gate{}
	reverseDeps := map[string][]string{}
	for _, g := range gates {
		gateByOut[g.out] = g
		reverseDeps[g.in1] = append(reverseDeps[g.in1], g.out)
		reverseDeps[g.in2] = append(reverseDeps[g.in2], g.out)
	}
	removingLastDep := func(from, dep string) bool {
		deps := reverseDeps[from]
		for i, d := range deps {
			if d == dep {
				reverseDeps[from] = append(deps[:i], deps[i+1:]...)
				return len(deps) == 1
			}
		}
		return false
	}
	sortedGates := make([]gate, 0, len(gates))
	placedGates := map[string]bool{}
	var tryPlacing func(out string)
	tryPlacing = func(out string) {
		if placedGates[out] || len(reverseDeps[out]) > 0 {
			return
		}
		g, ok := gateByOut[out]
		if !ok {
			return
		}
		sortedGates = append(sortedGates, g)
		placedGates[out] = true
		if removingLastDep(g.in1, g.out) {
			tryPlacing(g.in1)
		}
		if removingLastDep(g.in2, g.out) {
			tryPlacing(g.in2)
		}
	}

	for _, g := range gates {
		tryPlacing(g.out)
	}
	if len(sortedGates) != len(gates) {
		panic(fmt.Sprintf("not all gates are there: %v", sortedGates))
	}
	slices.Reverse(sortedGates)
	copy(gates, sortedGates)
	return
}

func (b board) output() (lines []string) {
	vals := map[string]int8{}
	inputs := make([]string, 0, len(b.inputs))
	for in, val := range b.inputs {
		inputs = append(inputs, in)
		vals[in] = val
	}
	sort.Strings(inputs)
	for _, in := range inputs {
		lines = append(lines, fmt.Sprintf("%s: %d", in, b.inputs[in]))
	}
	lines = append(lines, "")
	for _, g := range b.gates {
		var opName string
		switch g.op {
		case and:
			opName = "AND"
		case or:
			opName = "OR "
		case xor:
			opName = "XOR"
		}
		lines = append(lines, fmt.Sprintf("%s[%d] %s %s[%d] -> %s[%d]", g.in1, vals[g.in1], opName, g.in2, vals[g.in2], g.out, g.val))
		vals[g.out] = g.val
	}
	lines = append(lines, "")
	outputs := make([]string, 0, len(b.outputs))
	for out := range b.outputs {
		outputs = append(outputs, out)
	}
	sort.Strings(outputs)
	for _, out := range outputs {
		lines = append(lines, fmt.Sprintf("%s: %d", out, b.outputs[out]))
	}
	return
}

func (b board) outputGraphviz() (lines []string) {
	lines = append(lines, "digraph {")
	lines = append(lines, "  rankdir=LR;")
	gates := map[string]gate{}
	for _, g := range b.gates {
		gates[g.out] = g
	}
	label := func(name string) string {
		if val, ok := b.inputs[name]; ok {
			return fmt.Sprintf("\"%s:%d\"", name, val)
		}
		if g, ok := gates[name]; ok {
			var opName string
			switch g.op {
			case and:
				opName = "AND"
			case or:
				opName = "OR"
			case xor:
				opName = "XOR"
			}
			return fmt.Sprintf("\"%s:%s:%d\"", opName, name, g.val)
		}
		return "???"
	}
	for _, g := range b.gates {
		lines = append(lines, fmt.Sprintf("  %s -> %s;", label(g.in1), label(g.out)))
		lines = append(lines, fmt.Sprintf("  %s -> %s;", label(g.in2), label(g.out)))
	}
	lines = append(lines, "}")
	return
}

func (d Solve) CorrectAnswers() [2]string {
	return [2]string{"43942008931358", "dvb,fhg,fsq,tnc,vcf,z10,z17,z39"}
}
