package day09

import (
	"container/heap"
	"sort"
	"strings"
)

type Solve struct{}

func (d Solve) Part1(input []string, outputCh chan<- []string) int64 {
	blocks := parse(input)
	blocks = blocks.compactFragmenting()
	outputCh <- blocks.output()
	return blocks.checksum()
}

func (d Solve) Part2(input []string, outputCh chan<- []string) int64 {
	blocks := parse(input)
	blocks = blocks.compactContinuous()
	outputCh <- blocks.output()
	return blocks.checksum()
}

type block struct {
	id  int
	len byte
}

func (b block) isGap() bool {
	return b.id == -1
}

type blocks []block

type posBlock struct {
	block
	pos int
}

func parse(input []string) (blocks blocks) {
	pos := 0
	for _, line := range input {
		for id := 0; id*2 < len(line); id++ {
			blockLen := line[id*2] - '0'
			blocks = append(blocks, block{id: id, len: blockLen})
			pos += int(blockLen)
			if id*2+1 < len(line) {
				gapLen := line[id*2+1] - '0'
				blocks = append(blocks, block{id: -1, len: gapLen})
				pos += int(gapLen)
			}
		}
	}
	return
}

func (original blocks) compactFragmenting() (compacted blocks) {
	blocks := append(blocks{}, original...)
	i := 0
	j := len(blocks) - 1
	for i < j {
		if blocks[i].id >= 0 {
			compacted = append(compacted, blocks[i])
			i++
			continue
		}
		if blocks[j].id < 0 {
			j--
			continue
		}
		gap, blk := blocks[i], blocks[j]
		if gap.len < blk.len {
			compacted = append(compacted, block{id: blk.id, len: gap.len})
			blocks[j].len -= gap.len
			i++
			continue
		}
		if gap.len > blk.len {
			compacted = append(compacted, blk)
			blocks[i].len -= blk.len
			j--
			continue
		}
		compacted = append(compacted, blk)
		i++
		j--
	}
	if i == j && blocks[i].id >= 0 {
		compacted = append(compacted, blocks[i])
	}
	return
}

type gapHeap []*posBlock
type gapHeaps [10]gapHeap

func (h gapHeap) Len() int           { return len(h) }
func (h gapHeap) Less(i, j int) bool { return h[i].pos < h[j].pos }
func (h gapHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *gapHeap) Push(x any) {
	*h = append(*h, x.(*posBlock))
}
func (h *gapHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func (ghs *gapHeaps) push(gap posBlock) {
	for len := range gap.len {
		heap.Push(&ghs[len+1], &gap)
	}
}

func (ghs *gapHeaps) popFitting(posBlk posBlock) (foundGap posBlock, ok bool) {
nextLen:
	for i := int(posBlk.len); i < len(ghs); i++ {
		var gapPtr *posBlock
		for {
			if len(ghs[i]) == 0 {
				continue nextLen
			}
			gapPtr = heap.Pop(&ghs[i]).(*posBlock)
			if gapPtr.len > 0 {
				break
			}
		}
		gap := *gapPtr
		if gap.pos > posBlk.pos {
			continue nextLen
		}
		gapPtr.len = 0
		return gap, true
	}
	return posBlock{}, false
}

func (original blocks) compactContinuous() (compacted blocks) {
	gapHeaps := gapHeaps{}
	posOriginal := make([]posBlock, 0, len(original)/2+1)
	for i, pos := 0, 0; i < len(original); i, pos = i+1, pos+int(original[i].len) {
		posBlock := posBlock{original[i], pos}
		if posBlock.isGap() {
			gapHeaps.push(posBlock)
		} else {
			posOriginal = append(posOriginal, posBlock)
		}
	}
	posCompacted := make([]posBlock, 0, len(posOriginal))
	for i := len(posOriginal) - 1; i >= 0; i-- {
		posBlk := posOriginal[i]
		if posBlk.isGap() {
			continue
		}
		posGap, ok := gapHeaps.popFitting(posBlk)
		if !ok {
			posCompacted = append(posCompacted, posBlk)
			continue
		}
		posBlk.pos = posGap.pos
		posCompacted = append(posCompacted, posBlk)
		if posBlk.len == posGap.len {
			continue
		}
		posGap.len -= posBlk.len
		posGap.pos += int(posBlk.len)
		gapHeaps.push(posGap)
	}
	sort.Slice(posCompacted, func(i, j int) bool {
		return posCompacted[i].pos < posCompacted[j].pos
	})
	pos := 0
	for _, posBlk := range posCompacted {
		if pos < posBlk.pos {
			compacted = append(compacted, block{id: -1, len: byte(posBlk.pos - pos)})
		}
		compacted = append(compacted, posBlk.block)
		pos = posBlk.pos + int(posBlk.len)
	}
	return
}

func (blocks blocks) checksum() (sum int64) {
	i := 0
	for _, block := range blocks {
		if block.id < 0 {
			i += int(block.len)
			continue
		}
		for range block.len {
			sum += int64(block.id * i)
			i++
		}
	}
	return
}

func (blocks blocks) output() (lines []string) {
	maxLen := 64
	var sb strings.Builder
	for _, block := range blocks {
		if sb.Len()+int(block.len) > maxLen {
			lines = append(lines, sb.String())
			sb.Reset()
		}
		r := '.'
		if block.id >= 0 {
			r = '0' + rune(block.id%10)
		}
		sb.WriteString(strings.Repeat(string(r), int(block.len)))
	}
	if sb.Len() > 0 {
		lines = append(lines, sb.String())
	}
	return lines
}

func (d Solve) CorrectAnswers() [2]int64 {
	return [2]int64{6259790630969, 6289564433984}
}
