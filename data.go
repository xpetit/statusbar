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
	var received, transmitted int
	lines := readLines("/proc/net/dev")
	lines = lines[2:] // Ignore the 2-lines header
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 10 {
			continue // skip malformated lines
		}
		ifName := strings.TrimSuffix(fields[0], ":")
		switch ifName {
		case "docker0", "lo":
			continue // skip local network interfaces
		}
		rx, err := strconv.Atoi(fields[1])
		check(err)
		tx, err := strconv.Atoi(fields[9])
		check(err)
		received += rx
		transmitted += tx
	}
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
