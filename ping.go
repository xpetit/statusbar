package main

import (
	"fmt"
	"net"
	"sync/atomic"
	"time"
)

var (
	lastPing = make(chan time.Duration)
	lastTime atomic.Value
)

func init() {
	lastTime.Store(time.Now())
	go func() {
		for {
			t := time.Now()
			conn, err := net.Dial("tcp", "1.1.1.1:53")
			if err == nil {
				check(conn.Close())
			}
			lastTime.Store(time.Now())
			if err == nil {
				lastPing <- time.Since(t)
			} else {
				lastPing <- -1
			}
		}
	}()
}

// ping returns "XXXX ms", time to do a TCP connection on the Cloudflare DNS server
// or it returns " no net" if the connection failed or took more than 10 seconds
func ping() string {
	var ping time.Duration
	select {
	case ping = <-lastPing:
	default:
		ping = time.Since(lastTime.Load().(time.Time))
	}
	if ping < 0 || ping > 10*time.Second {
		return " no net"
	}
	return fmt.Sprintf("%4d ms", ping.Milliseconds())
}

// TODO: use true ICMP ping instead of TCP
