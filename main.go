package main

import (
	"fmt"
	"net"
	"os"
	"path"
	"runtime"
	"strconv"
	"time"

	"github.com/xpetit/statusbar/cpu"
	"github.com/xpetit/statusbar/mem"
	"github.com/xpetit/statusbar/network"
	"github.com/xpetit/statusbar/sensors"

	. "github.com/xpetit/x/v2"
)

func date() string {
	s := time.Now().Local().Format("Mon 02/01 15:04")
	return s[:2] + s[3:] // Remove the 3rd letter of the day since it is not differentiating
}

var cpuUsageWidth = len(strconv.Itoa(100 * runtime.NumCPU()))

func main() {
	socket := path.Join(os.Getenv("XDG_RUNTIME_DIR"), "statusbar.sock")
	os.Remove(socket)
	l := C2(net.Listen("unix", socket))
	for {
		conn := C2(l.Accept())
		fmt.Fprintf(conn, " %6s │ %.f° │ %*.f%% │ %.1fGB │ %s ",
			FormatByte(network.Usage()),
			sensors.MaxTemperature(),
			cpuUsageWidth, cpu.Usage(),
			float64(mem.BytesAvailable())/1e9,
			date(),
		)
		C(conn.Close())
	}
}
