package logs

import (
	"fmt"
	"sync/atomic"
	"syscall"
)

type severity int32 // sync/atomic int32

const (
	debugLog severity = iota
	infoLog
	warningLog
	errorLog
	fatalLog
	numSeverity = 5
)

var severityName = []string{
	debugLog:   "DEBUG",
	infoLog:    "INFO ",
	warningLog: "WARN ",
	errorLog:   "ERROR",
	fatalLog:   "FATAL",
}

// get returns the value of the severity.
func (s *severity) get() severity {
	return severity(atomic.LoadInt32((*int32)(s)))
}

// set sets the value of the severity.
func (s *severity) set(val severity) {
	atomic.StoreInt32((*int32)(s), int32(val))
}

var kernel32 = syscall.NewLazyDLL("kernel32.dll")

func brush(s string, i int) {
	proc := kernel32.NewProc("SetConsoleTextAttribute")
	handle, _, _ := proc.Call(uintptr(syscall.Stdout), uintptr(i))
	fmt.Printf("%s", s)
	handle, _, _ = proc.Call(uintptr(syscall.Stdout), uintptr(7))
	CloseHandle := kernel32.NewProc("CloseHandle")
	CloseHandle.Call(handle)
}
