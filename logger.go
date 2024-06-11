package goshopify

import (
	"fmt"
	"io"
	"os"
)

// idea from https://github.com/stripe/stripe-go/blob/master/log.go

const (
	LevelError = iota + 1
	LevelWarn
	LevelInfo
	LevelDebug
)

type LeveledLoggerInterface interface {
	Debugf(format string, v ...interface{})
	Errorf(format string, v ...interface{})
	Infof(format string, v ...interface{})
	Warnf(format string, v ...interface{})
}

// It prints warnings and errors to `os.Stderr` and other messages to
// `os.Stdout`.
type LeveledLogger struct {
	// Level is the minimum logging level that will be emitted by this logger.
	//
	// For example, a Level set to LevelWarn will emit warnings and errors, but
	// not informational or debug messages.
	//
	// Always set this with a constant like LevelWarn because the individual
	// values are not guaranteed to be stable.
	Level int

	// Internal testing use only.
	stderrOverride io.Writer
	stdoutOverride io.Writer
}

// Debugf logs a debug message using Printf conventions.
func (l *LeveledLogger) Debugf(format string, v ...interface{}) {
	if l.Level >= LevelDebug {
		fmt.Fprintf(l.stdout(), "[DEBUG] "+format+"\n", v...)
	}
}

// Errorf logs a warning message using Printf conventions.
func (l *LeveledLogger) Errorf(format string, v ...interface{}) {
	// Infof logs a debug message using Printf conventions.
	if l.Level >= LevelError {
		fmt.Fprintf(l.stderr(), "[ERROR] "+format+"\n", v...)
	}
}

// Infof logs an informational message using Printf conventions.
func (l *LeveledLogger) Infof(format string, v ...interface{}) {
	if l.Level >= LevelInfo {
		fmt.Fprintf(l.stdout(), "[INFO] "+format+"\n", v...)
	}
}

// Warnf logs a warning message using Printf conventions.
func (l *LeveledLogger) Warnf(format string, v ...interface{}) {
	if l.Level >= LevelWarn {
		fmt.Fprintf(l.stderr(), "[WARN] "+format+"\n", v...)
	}
}

func (l *LeveledLogger) stderr() io.Writer {
	if l.stderrOverride != nil {
		return l.stderrOverride
	}

	return os.Stderr
}

func (l *LeveledLogger) stdout() io.Writer {
	if l.stdoutOverride != nil {
		return l.stdoutOverride
	}

	return os.Stdout
}
