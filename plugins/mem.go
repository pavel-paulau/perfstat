package plugins

import (
	"bufio"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

type Mem struct {
	Columns []string
}

func NewMem() *Mem {
	columns := []string{
		"mem_used",  // Total RAM in use (doesn't include buffers and cache).
		"mem_free",  // Amount of free RAM.
		"mem_buff",  // Relatively temporary storage for raw disk blocks.
		"mem_cache", // In-memory cache for files read from the disk.
	}

	return &Mem{
		Columns: columns,
	}
}

func (m *Mem) GetColumns() []string {
	return m.Columns
}

func (m *Mem) Extract() (results []float64) {
	file, err := os.Open("/proc/meminfo")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	for i := 0; i < len(m.Columns); i++ {
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Println(err)
		}

		value, err := strconv.ParseFloat(strings.Fields(line)[1], 32)
		value = math.Floor(value / 1024)
		if err != nil {
			log.Println(err)
		}
		results = append(results, value)
	}
	for i := 1; i < len(results); i++ {
		results[0] -= results[i]
	}
	return
}
