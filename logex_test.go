package logex

import (
	"testing"
)

type S struct {
	Logger
}

func (s *S) hello() {
	s.Warn("warn in hello")
}

var s = S{}

func TestLogex(t *testing.T) {
	log("a")
	Info("b")
	NewLogger(0).Info("c")
	Error("ec")
	s.hello()
	Struct(&s, 1, "", false)
}

func log(s string) {
	Logger{}.Output(2, s)
}
