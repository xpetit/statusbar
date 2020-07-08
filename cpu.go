package main

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var (
	threads = runtime.NumCPU()
	width   = len(strconv.Itoa(100 * threads))
)

type cpu struct {
	idle int
	time time.Time
}

// returns "XXX %" where "100 %" means one thread is busy
func (last *cpu) String() string {
	lines := readLines("/proc/stat")
	fields := strings.Fields(lines[0])
	idle, err := strconv.Atoi(fields[4])
	check(err)
	now := time.Now()
	delta := float64(now.Sub(last.time)) / 1_000_000_000
	percent := 100 * (float64(threads) - (float64(idle-last.idle) / delta / 100))
	last.idle = idle
	last.time = now
	return fmt.Sprintf("%*.f %%", width, percent)[:width+2]
}
