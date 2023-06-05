package cpu

import (
	"bufio"
	"io"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	. "github.com/xpetit/x/v2"
)

var (
	stat     = C2(os.Open("/proc/stat"))
	ncpus    = float64(runtime.NumCPU())
	lastTime time.Time
	lastIdle float64
)

// Usage returns the CPU thread utilization in percent since last call (100 % means one thread is busy)
func Usage() float64 {
	C2(stat.Seek(0, io.SeekStart))
	var (
		line    = C2(bufio.NewReader(stat).ReadString('\n'))
		idle    = float64(C2(strconv.Atoi(strings.Fields(line)[4])))
		now     = time.Now()
		secs    = now.Sub(lastTime).Seconds()
		percent = 100 * (ncpus - (idle-lastIdle)/secs/100)
	)
	lastIdle = idle
	lastTime = now
	return percent
}
