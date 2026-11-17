//go:build !linux

package sandbox

import (
	"fmt"
	"syscall"
)

// Detect: no sandbox but on Linux.
func Detect(m Mode) (Mode, error) {
	if m == ModeOff {
		return m, nil
	}
	return "", fmt.Errorf("%w: only on Linux", ErrUnavailable)
}

// LandlockABI: none.
func LandlockABI() int { return 0 }

// Proc is a program started.
type Proc struct{ Pid int }

func (p *Proc) Wait() error               { return nil }
func (p *Proc) Done() <-chan struct{}     { c := make(chan struct{}); close(c); return c }
func (p *Proc) Signal(sig syscall.Signal) {}

// Start: not on this system.
func Start(m Mode, s Spec) (*Proc, error) {
	return nil, fmt.Errorf("%w: only on Linux", ErrUnavailable)
}

// CoverageOf: nothing kept apart but on Linux.
func CoverageOf(m Mode) Coverage { return Coverage{Mode: m} }
