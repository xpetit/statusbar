package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

const filename = "/tmp/statusbar"

func check(err error) {
	if err != nil {
		panic(err)
	}
}
func readLines(filename string) []string {
	b, err := ioutil.ReadFile(filename)
	check(err)
	return strings.Split(string(b), "\n")
}
func main() {
	data := new(data)
	cpu := new(cpu)
	for range time.Tick(3 * time.Second) {
		s := fmt.Sprintln(data, ping(), cpu, mem(), date())
		check(ioutil.WriteFile(filename+"_new", []byte(s), os.ModePerm))
		check(os.Rename(filename+"_new", filename))
	}
}

// TODO: add temperature
