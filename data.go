package main

import (
	"fmt"
	"math"
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
	// Collect routed network interfaces with a gateway
	activeInterfaces := map[string]struct{}{}
	{
		lines := readLines("/proc/net/route")
		lines = lines[1:] // Skip the header
		for _, line := range lines {
			fields := strings.Fields(line)
			if len(fields) < 4 {
				continue
			}
			flags, err := strconv.Atoi(fields[3])
			check(err)
			const (
				up      = 1 << iota // route usable
				gateway             // destination is a gateway
			)
			if flags&up != 0 && flags&gateway != 0 {
				ifName := fields[0]
				activeInterfaces[ifName] = struct{}{}
			}
		}
	}

	// Gather received and transmitted bytes for all routed network interfaces
	var received, transmitted int
	{
		lines := readLines("/proc/net/dev")
		lines = lines[2:] // Skip the 2-lines header
		for _, line := range lines {
			fields := strings.Fields(line)
			if len(fields) < 10 {
				continue // skip malformated lines
			}
			ifName := strings.TrimSuffix(fields[0], ":")
			if _, ok := activeInterfaces[ifName]; !ok {
				continue // skip network interfaces without a gateway
			}
			rx, err := strconv.Atoi(fields[1])
			check(err)
			tx, err := strconv.Atoi(fields[9])
			check(err)
			received += rx
			transmitted += tx
		}
	}

	// Compute and format stats
	now := time.Now()
	delta := float64(now.Sub(last.time)) / 1_000_000_000
	receivedPerSec := int(math.Round(float64(received-last.received) / delta))
	transmittedPerSec := int(math.Round(float64(transmitted-last.transmitted) / delta))
	last.received = received
	last.transmitted = transmitted
	last.time = now
	return fmt.Sprintf("%6s %6s", prettyByte(receivedPerSec), prettyByte(transmittedPerSec))
}

func prettyByte(i int) string {
	f := float64(i)
	if f < 0 {
		f = 0
	}
	unit := " kMGTPEZY"
	for f >= 1000 {
		unit = unit[1:]
		f /= 1000
	}
	var prec int
	if unit[0] != ' ' && math.Round(f) < 10 {
		prec = 1
	}
	return fmt.Sprintf("%3.*f %cB", prec, f, unit[0])
}
