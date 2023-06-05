package mem

import (
	"bufio"
	"io"
	"os"
	"strconv"
	"strings"

	. "github.com/xpetit/x/v2"
)

var f = C2(os.Open("/proc/meminfo"))

// BytesAvailable return the amount of free RAM (in bytes)
func BytesAvailable() int {
	C2(f.Seek(0, io.SeekStart))
	scanner := bufio.NewScanner(f)
	const key = "MemAvailable:"
	for scanner.Scan() {
		s := scanner.Text()
		if !strings.HasPrefix(s, key) {
			continue
		}
		s = s[len(key):]              // skip key
		s = strings.TrimLeft(s, " ")  // skip spaces
		s, _, _ = strings.Cut(s, " ") // cut before spaces
		kib := C2(strconv.Atoi(s))
		return kib * 1024 // KiB to B
	}
	C(scanner.Err())
	panic(key + " not found in /proc/meminfo")
}
