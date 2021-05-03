/*
 * Copyright (c) 2019-present Heeus authors
 */

package service

import "fmt"

const ErrLevel = 0
const DebugLevel = 1

type Logger struct {
	level int64
	//showTS bool
}

func (l *Logger) Error(msg string, args ...interface{}) {
	if l.level < ErrLevel {
		return
	}

	fmt.Printf("Error: "+msg+"\n", args...)
}

func (l *Logger) Debug(msg string, args ...interface{}) {
	if l.level < DebugLevel {
		return
	}

	fmt.Printf("Debug: "+msg+"\n", args...)
}

func (l *Logger) Log(msg string, args ...interface{}) {
	fmt.Printf(msg+"\n", args...)
}
