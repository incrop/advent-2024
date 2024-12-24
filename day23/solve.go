package day23

import (
	"incrop/advent-2024/out"
	"sort"
	"strings"
)

type Solve struct{}

func (d Solve) Part1(input []string, outputCh chan<- []string) (count int64) {
	net := parse(input)
	l := out.NewLog(outputCh)
	for _, group := range net.fullyConnectedGroups() {
		if len(group) != 3 {
			continue
		}
		a, b, c := group[0], group[1], group[2]
		na, nb, nc := net.nodes[a], net.nodes[b], net.nodes[c]
		if strings.HasPrefix(na, "t") || strings.HasPrefix(nb, "t") || strings.HasPrefix(nc, "t") {
			l.Printf("%s,%s,%s", na, nb, nc)
			count++
		}
	}
	return
}

func (d Solve) Part2(input []string, outputCh chan<- []string) int64 {
	net := parse(input)
	l := out.NewLog(outputCh)
	var maxGroup []int
	for _, group := range net.fullyConnectedGroups() {
		if len(group) > len(maxGroup) {
			maxGroup = group
		}
	}
	maxGroupNodes := make([]string, len(maxGroup))
	for i, n := range maxGroup {
		maxGroupNodes[i] = net.nodes[n]
	}
	l.Printf("%s", strings.Join(maxGroupNodes, ","))
	l.Append(net.output())
	return 0
}

type network struct {
	nodes []string
	conn  [][]bool
}

func (net network) fullyConnectedGroups() (groups [][]int) {
	groups = [][]int{}
	for i := range net.nodes {
		next := groups
		for _, group := range groups {
			allConnected := true
			for _, j := range group {
				if !net.conn[j][i-j-1] {
					allConnected = false
					break
				}
			}
			if allConnected {
				biggerGroup := make([]int, len(group)+1)
				copy(biggerGroup, group)
				biggerGroup[len(group)] = i
				next = append(next, biggerGroup)
			}
		}
		next = append(next, []int{i})
		groups = next
	}
	return
}

func parse(input []string) (net network) {
	links := [][2]string{}
	nodes := map[string]bool{}
	for _, line := range input {
		link := strings.Split(line, "-")
		if link[0] > link[1] {
			link[0], link[1] = link[1], link[0]
		}
		links = append(links, [2]string{link[0], link[1]})
		nodes[link[0]], nodes[link[1]] = true, true
	}
	net.nodes = make([]string, 0, len(nodes))
	for node := range nodes {
		net.nodes = append(net.nodes, node)
	}
	sort.Strings(net.nodes)
	nodeIdx := map[string]int{}
	net.conn = make([][]bool, len(net.nodes))
	for i, node := range net.nodes {
		nodeIdx[node] = i
		net.conn[i] = make([]bool, len(net.nodes)-i-1)
	}
	for _, link := range links {
		from, to := nodeIdx[link[0]], nodeIdx[link[1]]
		net.conn[from][to-from-1] = true
	}
	return
}

func (net network) output() (lines []string) {
	for i, node := range net.nodes {
		var sb strings.Builder
		sb.WriteString(node)
		sb.WriteString(" ->")
		for j, connected := range net.conn[i] {
			if connected {
				sb.WriteString(" ")
				sb.WriteString(net.nodes[i+j+1])
			}
		}
		lines = append(lines, sb.String())
	}
	return
}

func (d Solve) CorrectAnswers() [2]int64 {
	return [2]int64{1054, 0}
}
