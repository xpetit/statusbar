package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type data struct {
	received, transmitted int
	time                  time.Time
}

// returns "XXX [U]B YYY [U]B" where XXX and YYY are numbers from 0 to 999, [U] is a SI prefix
// XXX is the number of bytes received, YYY the number of bytes transmitted
func (last *data) String() string {
	// Collect network interfaces with a route
	activeInterfaces := map[string]struct{}{}
	routeLines := readLines("/proc/net/route")
	routeLines = routeLines[1:] // Skip the header
	for _, line := range routeLines {
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}
		activeInterfaces[fields[0]] = struct{}{}
	}

	// Gather received and transmitted bytes for all routed network interfaces
	var received, transmitted int
	devLines := readLines("/proc/net/dev")
	devLines = devLines[2:] // Skip the 2-lines header
	for _, line := range devLines {
		fields := strings.Fields(line)
		if len(fields) < 10 {
			continue // skip malformated lines
		}
		ifName := strings.TrimSuffix(fields[0], ":")
		if _, ok := activeInterfaces[ifName]; !ok {
			continue // skip unrouted network interfaces
		}
		rx, err := strconv.Atoi(fields[1])
		check(err)
		tx, err := strconv.Atoi(fields[9])
		check(err)
		received += rx
		transmitted += tx
	}

	// Compute and format stats
	now := time.Now()
	delta := float64(now.Sub(last.time)) / 1_000_000_000
	receivedPerSec := int(float64(received-last.received) / delta)
	transmittedPerSec := int(float64(transmitted-last.transmitted) / delta)
	last.received = received
	last.transmitted = transmitted
	last.time = now
	return fmt.Sprintf("%s %s", byteCountDecimal(receivedPerSec), byteCountDecimal(transmittedPerSec))
}

// From https://programming.guide/go/formatting-byte-size-to-human-readable-format.html
func byteCountDecimal(b int) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%3d  B", b)
	}
	div := unit
	exp := 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%3.f %cB", float64(b)/float64(div), "kMGTPEZY"[exp])
}

// TODO: consider changing value approximation:
//   123  B => 0.1 kB
//   123 kB => 0.1 MB
//   789 MB => 0.8 GB
