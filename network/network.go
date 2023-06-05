package network

import (
	"os"
	"strconv"
	"strings"
	"time"

	. "github.com/xpetit/x/v2"
)

func readLines(filename string) []string {
	return strings.Split(strings.TrimSpace(string(C2(os.ReadFile(filename)))), "\n")
}

var (
	lastTime      time.Time
	lastExchanged float64
)

// Usage returns the number of bytes exchanged (received & sent) since last call
func Usage() (exchangedPerSec float64) {
	// Collect routed network interfaces with a gateway
	activeInterfaces := map[string]struct{}{}
	{
		lines := readLines("/proc/net/route")[1:] // Skip the header
		for _, line := range lines {
			fields := strings.SplitN(line, "\t", 5)
			if len(fields) < 4 {
				continue
			}
			flags := C2(strconv.Atoi(fields[3]))
			const (
				up      = 1 << iota // route usable
				gateway             // destination is a gateway
			)
			if flags&(up|gateway) != 0 {
				ifName := fields[0]
				activeInterfaces[ifName] = struct{}{}
			}
		}
	}

	// Gather received and transmitted bytes for all routed network interfaces
	var exchanged float64
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
			rx := float64(C2(strconv.Atoi(fields[1])))
			tx := float64(C2(strconv.Atoi(fields[9])))
			exchanged += rx
			exchanged += tx
		}
	}
	now := time.Now()
	exchangedPerSec = (exchanged - lastExchanged) / now.Sub(lastTime).Seconds()
	lastExchanged = exchanged
	lastTime = now
	return exchangedPerSec
}
