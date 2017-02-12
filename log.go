package seared

import (
	"log"
	"testing"
)

type Log interface {
	Debugf(format string, args ...interface{})
}

func StandardLog() Log {
	return &standardLog{}
}

type standardLog struct{}

func (l *standardLog) Debugf(format string, args ...interface{}) {
	log.Printf(format, args...)
}

func TestingLog(t *testing.T) Log {
	return &testingLog{t}
}

type testingLog struct {
	t *testing.T
}

func (l *testingLog) Debugf(format string, args ...interface{}) {
	l.t.Logf(format, args...)
}
