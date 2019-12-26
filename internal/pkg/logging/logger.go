// Copyright 2019 The JIBit Team. All rights reserved.
// Use of this source code is governed by an Apache Style license that can be found in the LICENSE file.

// Package logging provides utility methods based on standard log.Logger implementation.
// This package supports different log levels with fluent interfaces.
package logging

import (
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"strings"
)

const (
	// Defines log pattern; Including filename and line number of logger caller and log message.
	pattern = "[ %s %d ] -- %s"

	// Defines flags pattern that is related to date and time of a log.
	flags = log.Ldate | log.Ltime | log.Lmicroseconds
)

// Level defines log level type, including DEBUG, INFO, WARN and ERROR.
type Level int

// String is the implementation of Stringer interface for level.
func (l Level) String() string {
	return []string{"[ DEBUG ] ", "[ INFO ] ", "[ WARN ] ", "[ ERROR ] "}[l]
}

// Different log levels definition.
const (
	DEBUG Level = iota
	INFO
	WARN
	ERROR
)

// Logger defines a logger implementation with support for different log levels.
type Logger struct {
	level Level
	d     *log.Logger // Debug level logger.
	i     *log.Logger // Info level logger.
	w     *log.Logger // Warning level logger.
	e     *log.Logger // Error level logger.
}

// New creates a new logger with INFO as log level and standard error stream as output.
func New() *Logger {
	return &Logger{
		level: INFO,
		d:     log.New(os.Stderr, DEBUG.String(), flags),
		i:     log.New(os.Stderr, INFO.String(), flags),
		w:     log.New(os.Stderr, WARN.String(), flags),
		e:     log.New(os.Stderr, ERROR.String(), flags),
	}
}

// NewWithLevel creates a new logger with specified log level as string value.
func NewWithLevel(level string) *Logger {
	switch strings.ToLower(level) {
	case "debug":
		return New().WithLevel(DEBUG)
	case "info":
		return New().WithLevel(INFO)
	case "warn", "warning":
		return New().WithLevel(WARN)
	case "error":
		return New().WithLevel(ERROR)
	default:
		return New().WithLevel(INFO)
	}
}

// WithLevel updates level of logger.
func (l *Logger) WithLevel(level Level) *Logger {
	l.level = level

	return l
}

// WithWriter updates writer (io.Write) of logger. For example it is useful if you want to log on a file.
func (l *Logger) WithWriter(w io.Writer) *Logger {
	l.d = log.New(w, DEBUG.String(), flags)
	l.i = log.New(w, INFO.String(), flags)
	l.w = log.New(w, WARN.String(), flags)
	l.e = log.New(w, ERROR.String(), flags)

	return l
}

// Debug prints a debug logging message if log level is equal to debug.
func (l *Logger) Debug(message string, args ...interface{}) {
	if l.level == DEBUG {
		file, line := caller()
		l.d.Printf(pattern, file, line, fmt.Sprintf(message, args...))
	}
}

// Info prints an info logging message if log level is less than or equal to info.
func (l *Logger) Info(message string, args ...interface{}) {
	if l.level <= INFO {
		file, line := caller()
		l.i.Printf(pattern, file, line, fmt.Sprintf(message, args...))
	}
}

// Warn prints a warning level logging message if log level is less than or equal to warning.
func (l *Logger) Warn(message string, args ...interface{}) {
	if l.level <= WARN {
		file, line := caller()
		l.w.Printf(pattern, file, line, fmt.Sprintf(message, args...))
	}
}

// Error prints an error logging message if log level is less than or equal to error.
func (l *Logger) Error(message string, args ...interface{}) {
	if l.level <= ERROR {
		file, line := caller()
		l.e.Printf(pattern, file, line, fmt.Sprintf(message, args...))
	}
}

// Fatal prints an error logging message if log level is less than or equal to error followed by a call to os.Exit(1).
func (l *Logger) Fatal(message string, args ...interface{}) {
	file, line := caller()
	l.e.Fatalf(pattern, file, line, fmt.Sprintf(message, args...))
}

// Extracts filename and line number of the logger caller. See runtime.Caller() for more information.
func caller() (string, int) {
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "????"
		line = 0
		return file, line
	}

	return file[1:], line
}
