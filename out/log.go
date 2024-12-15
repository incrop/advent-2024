package out

import (
	"fmt"
	"time"
)

type Log struct {
	ch    chan<- []string
	lines []string
	delay time.Duration
}

func NewLog(ch chan<- []string) *Log {
	return &Log{
		ch:    ch,
		lines: nil,
	}
}

func (l *Log) WithDelay(delay time.Duration) *Log {
	l.delay = delay
	return l
}

func (l *Log) Printf(format string, args ...any) {
	line := fmt.Sprintf(format, args...)
	l.lines = append(l.lines, line)
	l.ch <- l.lines
	if l.delay > 0 {
		time.Sleep(l.delay)
	}
}

func (l *Log) Append(lines []string) {
	l.lines = append(l.lines, lines...)
	l.ch <- l.lines
	if l.delay > 0 {
		time.Sleep(l.delay)
	}
}
