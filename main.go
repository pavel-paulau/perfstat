package main

import (
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	"code.google.com/p/go.crypto/ssh/terminal"

	"github.com/pavel-paulau/perfstat/plugins"
)

var header []string
var values []float64

func printHeader() {
	hr := 0
	for _, column := range header {
		fmt.Printf("%s ", column)
		hr += len(column) + 1
	}
	fmt.Println()
	fmt.Println(strings.Repeat("-", hr-1))
}

func printValues() {
	for i, column := range header {
		fmt_str := fmt.Sprintf("%%%dv ", len(column))
		fmt.Printf(fmt_str, values[i])
	}
	fmt.Println()
	values = []float64{}
}

func main() {
	cpu := flag.Bool("cpu", false, "enable CPU stats")
	mem := flag.Bool("mem", false, "enable memory stats")
	interval := flag.Int("interval", 1, "sampling interval in seconds")
	flag.Parse()

	active_plugins := []plugins.Plugin{}
	if *cpu == true {
		active_plugins = append(active_plugins, plugins.NewCPU())
	}
	if *mem == true {
		active_plugins = append(active_plugins, plugins.NewMem())
	}

	if len(active_plugins) == 0 {
		log.Fatalln("Please specify at least one plugin")
	}

	for _, plugin := range active_plugins {
		header = append(header, plugin.GetColumns()...)
	}
	printHeader()

	_, y, err := terminal.GetSize(0)
	if err != nil {
		log.Fatalln(err)
	}

	iterations := 1
	for {
		for _, plugin := range active_plugins {
			values = append(values, plugin.Extract()...)
		}
		printValues()

		iterations += 1
		if iterations == y-1 {
			printHeader()
			iterations = 1
		}
		time.Sleep(time.Duration(*interval) * time.Second)
	}
}
