package out

import (
	"fmt"
)

type Log struct {
	ch    chan<- []string
	lines []string
}

func NewLog(ch chan<- []string) *Log {
	return &Log{
		ch:    ch,
		lines: nil,
	}
}

func (l *Log) Printf(format string, args ...any) {
	line := fmt.Sprintf(format, args...)
	l.lines = append(l.lines, line)
	l.ch <- l.lines
}
