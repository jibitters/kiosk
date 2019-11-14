// Copyright 2019 The Jibit Team. All rights reserved.
// Use of this source code is governed by a Apache-style license that can be found in the LICENSE.md file.

// Package logging contains some utility functions related to logging.
package logging

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
)

// Defines log pattern; Including filename and line number of logger caller, and log message.
const logPattern = "[ %s %d ] -- %s"

// Defines log level type.
type level int

// Different log level constants.
const (
	DebugLevel level = iota
	InfoLevel
	WarningLevel
	ErrorLevel
)

// Logger defines a logger implementation with support for different log levels.
// Logger uses different go standard loggers for each log level, to reduce contention.
type Logger struct {
	debugLogger   *log.Logger
	infoLogger    *log.Logger
	warningLogger *log.Logger
	errorLogger   *log.Logger
	level         level
}

// New creates a new logger.
func New(level level) *Logger {
	return &Logger{
		debugLogger:   log.New(os.Stderr, "[ DEBUG ] ", log.LstdFlags),
		infoLogger:    log.New(os.Stderr, "[ INFO  ] ", log.LstdFlags),
		warningLogger: log.New(os.Stderr, "[ WARN  ] ", log.LstdFlags),
		errorLogger:   log.New(os.Stderr, "[ ERROR ] ", log.LstdFlags),
		level:         level,
	}
}

// NewWithLevel creates a new logger with specified log level as string value.
func NewWithLevel(level string) *Logger {
	switch strings.ToLower(level) {
	case "debug":
		return New(DebugLevel)
	case "info":
		return New(InfoLevel)
	case "warn", "warning":
		return New(WarningLevel)
	case "error":
		return New(ErrorLevel)
	}

	return New(InfoLevel)
}

// Debug prints a debug logging message if log level is equal to debug.
func (logger *Logger) Debug(message string, args ...interface{}) {
	if logger.level == DebugLevel {
		file, line := callerInfo()
		logger.debugLogger.Printf(logPattern, file, line, fmt.Sprintf(message, args...))
	}
}

// Info prints an info logging message if log level is less than or equal to info.
func (logger *Logger) Info(message string, args ...interface{}) {
	if logger.level <= InfoLevel {
		file, line := callerInfo()
		logger.infoLogger.Printf(logPattern, file, line, fmt.Sprintf(message, args...))
	}
}

// Warning prints a warning level logging message if log level is less than or equal to warning.
func (logger *Logger) Warning(message string, args ...interface{}) {
	if logger.level <= WarningLevel {
		file, line := callerInfo()
		logger.warningLogger.Printf(logPattern, file, line, fmt.Sprintf(message, args...))
	}
}

// Error prints an error logging message if log level is less than or equal to error.
func (logger *Logger) Error(message string, args ...interface{}) {
	if logger.level <= ErrorLevel {
		file, line := callerInfo()
		logger.errorLogger.Printf(logPattern, file, line, fmt.Sprintf(message, args...))
	}
}

// Fatal prints an error logging message if log level is less than or equal to error followed by a call to os.Exit(1).
func (logger *Logger) Fatal(message string, args ...interface{}) {
	if logger.level <= ErrorLevel {
		file, line := callerInfo()
		logger.errorLogger.Fatalf(logPattern, file, line, fmt.Sprintf(message, args...))
	}
}

// Extracts filename and line number of the logger caller. See runtime.Caller() for more information.
func callerInfo() (string, int) {
	var _, file, line, ok = runtime.Caller(2)
	if !ok {
		file = "???"
		line = 0
		return file, line
	}

	return shortFile(file), line
}

// Tries to short the complete relative file path and file name to only file name.
func shortFile(file string) string {
	short := file
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			short = file[i+1:]
			break
		}
	}

	return short
}
