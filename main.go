package main

import (
	"fmt"
	"log"
	"time"

	"code.google.com/p/go.crypto/ssh/terminal"

	"github.com/pavel-paulau/perfstat/plugins"
)

var header []string
var values []float64

var INTERVAL = 1 * time.Second

func printHeader() {
	for _, column := range header {
		fmt.Printf("%s ", column)
	}
	fmt.Println()
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
	_, y, err := terminal.GetSize(0)
	if err != nil {
		log.Fatalln(err)
	}

	plugin := plugins.NewCPU()
	header = append(header, plugin.Columns...)
	printHeader()

	iterations := 1
	for {
		values = append(values, plugin.Extract()...)
		printValues()

		iterations += 1
		if iterations == y {
			printHeader()
			iterations = 1
		}
		time.Sleep(INTERVAL)
	}
}
