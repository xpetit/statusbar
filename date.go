package main

import "time"

func date() string {
	now := time.Now().Local()
	day := now.Format("Mon")[:2]
	return now.Format(day + " 02/01 15:04")
}
