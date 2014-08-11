package plugins

import (
	"bufio"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

type CPU struct {
	previous, current []float64
	Columns           []string
}

func NewCPU() *CPU {
	columns := []string{
		"cpu_user",
		"cpu_nice",
		"cpu_sys",
		"cpu_idle",
		"cpu_iowait",
	}

	return &CPU{
		previous: make([]float64, len(columns)+1),
		Columns:  columns,
	}
}

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
	stats := strings.Fields(line)[1:] // omit "cpu" column

	total := float64(0)
	for i, _ := range c.Columns {
		value, err := strconv.ParseFloat(stats[i], 32)
		if err != nil {
			log.Println(err)
		}
		c.current = append(c.current, value)
		total += value
	}
	c.current = append(c.current, total)

	total_time := c.current[len(c.current)-1] - c.previous[len(c.previous)-1]
	for i, _ := range c.Columns {
		results = append(results, math.Floor(100*(c.current[i]-c.previous[i])/total_time+0.5))
	}
	c.previous = c.current
	c.current = []float64{}
	return
}
