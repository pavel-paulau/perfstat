package main

import (
	"flag"
	"fmt"
	"os"
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
		fmtStr := fmt.Sprintf("%%%dv ", len(column))
		fmt.Printf(fmtStr, values[i])
	}
	fmt.Println()
	values = []float64{}
}

func main() {
	cpu := flag.Bool("cpu", false, "enable CPU stats (Irix mode)")
	mem := flag.Bool("mem", false, "enable memory stats")

	interval := flag.Int("interval", 1, "sampling interval in seconds")

	quiet := flag.Bool("quiet", false, "disable reporting to stdout")
	perfkeeper := flag.String("perfkeeper", "127.0.0.1:8080", "optional perfkeeper host:port")
	snapshot := flag.String("snapshot", "", "name of perfkeeper snapshot")
	source := flag.String("source", "", "name of perfkeeper snapshot")

	flag.Parse()

	activePlugins := []plugins.Plugin{}
	if *cpu == true {
		activePlugins = append(activePlugins, plugins.NewCPU())
	}
	if *mem == true {
		activePlugins = append(activePlugins, plugins.NewMem())
	}
	if len(activePlugins) == 0 {
		fmt.Println("Please specify at least one plugin")
		os.Exit(1)
	}

	var keeper *keeper
	if *snapshot != "" && *source != "" {
		keeper = newKeeper(*perfkeeper, *snapshot, *source)
	} else {
		keeper = nil
	}

	for _, plugin := range activePlugins {
		header = append(header, plugin.GetColumns()...)
	}
	if !*quiet {
		printHeader()
	}

	_, terminalHeight, err := terminal.GetSize(0)
	if err != nil {
		terminalHeight = -1 // Detached terminal
	}

	iterations := 1
	for {
		for _, plugin := range activePlugins {
			values = append(values, plugin.Extract()...)
		}
		if keeper != nil {
			go keeper.store(header, values)
		}
		if !*quiet {
			printValues()
		}

		iterations++
		if iterations == terminalHeight-1 {
			if !*quiet {
				printHeader()
			}
			iterations = 1
		}
		time.Sleep(time.Duration(*interval) * time.Second)
	}
}
