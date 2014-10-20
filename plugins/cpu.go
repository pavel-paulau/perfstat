package plugins

import (
	"bufio"
	"log"
	"math"
	"os"
	"runtime"
	"strconv"
	"strings"
)

type CPU struct {
	Columns           []string
	nproc             float64
	previous, current []float64
}

// NewCPU initialize CPU plugin
func NewCPU() *CPU {
	columns := []string{
		"cpu_user",   // Time spent in user mode.
		"cpu_sys",    // Time spent in system mode.
		"cpu_idle",   // Time spent in the idle task.
		"cpu_iowait", // Time waiting for I/O to complete.
	}

	return &CPU{
		Columns:  columns,
		nproc:    float64(runtime.NumCPU()),
		previous: make([]float64, len(columns)+1),
	}
}

func (c *CPU) GetColumns() []string {
	return c.Columns
}

// Extract measures CPU utilization based on kernel/system statistics:
//	1. Take the first line of /proc/stat, e.g.:
//		cpu  1198438 862 391327 16130978 10149 53 14464 0 0 0)
//	2. Take the following values:
//		[1] user
//		[3] system
//		[4] idle
//		[5] iowait
// See `man proc` for details.
func (c *CPU) Extract() (results []float64) {
	file, err := os.Open("/proc/stat")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	line, err := reader.ReadString('\n')
	if err != nil {
		log.Println(err)
	}
	stats := append(strings.Fields(line)[1:2], strings.Fields(line)[3:6]...)

	total := float64(0)
	for i := range c.Columns {
		value, err := strconv.ParseFloat(stats[i], 32)
		if err != nil {
			log.Println(err)
		}
		c.current = append(c.current, value)
		total += value
	}
	c.current = append(c.current, total)

	totalTime := c.current[len(c.current)-1] - c.previous[len(c.previous)-1]
	for i := range c.Columns {
		results = append(results, math.Floor(c.nproc*100*(c.current[i]-c.previous[i])/totalTime+0.5))
	}
	c.previous = c.current
	c.current = []float64{}
	return
}
