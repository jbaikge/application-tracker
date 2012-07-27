package main

import (
	"bytes"
	"fmt"
	"os"
)

type App struct {
	Name string
	PID  int
}

var (
	initOffset int64 = 6
	maxLen     int64 = 15
)

func (a App) ProcessName() (name string) {
	f, err := os.Open(fmt.Sprintf("/proc/%d/status", a.PID))
	if err != nil {
		name = "N/A"
		return
	}
	// /proc/<PID>/stat[us] only contain the first 15 characters of the
	// executable's basename
	//
	// In /proc/<PID>/status, the name is followed by "Name:\t"
	b := make([]byte, 1)
	buf := bytes.NewBuffer(make([]byte, maxLen))
	for off := initOffset; off-initOffset < maxLen; off++ {
		if _, err := f.ReadAt(b, off); err != nil {
			break
		}
		if b[0] == '\n' {
			break
		}
		buf.WriteByte(b[0])
	}
	name = buf.String()
	return
}
