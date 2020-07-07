package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

var total = readMem("MemTotal:")

func readMem(parameter string) int {
	lines := readLines("/proc/meminfo")
	for _, line := range lines {
		fields := strings.Fields(line)
		if fields[0] == parameter {
			i, err := strconv.Atoi(fields[1])
			check(err)
			return int(math.Round(float64(i*1024) / 1_000_000_000)) // Kio to GB
		}
	}
	panic(parameter + " not found in /proc/meminfo")
}

// mem returns "used+available GB"
func mem() string {
	available := readMem("MemAvailable:") - 1 // minus 1 GB because Linux...
	used := total - available
	padding := 2*len(strconv.Itoa(total)) + 1
	return fmt.Sprintf("%*s GB", padding, fmt.Sprintf("%d+%d", used, available))
}
