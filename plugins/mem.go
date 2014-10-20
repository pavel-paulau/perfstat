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

// NewMem initialize Mem plugin
func NewMem() *Mem {
	columns := []string{
		"mem_used",  // Total RAM in use (doesn't include buffers and cache), Mbytes.
		"mem_free",  // Amount of free RAM, Mbytes.
		"mem_buff",  // Relatively temporary storage for raw disk blocks, Mbytes.
		"mem_cache", // In-memory cache for files read from the disk, Mbytes.
	}

	return &Mem{
		Columns: columns,
	}
}

func (m *Mem) GetColumns() []string {
	return m.Columns
}

// Extract measures RAM usage based kernel/system statistics:
//	1. Read 4 first lines from /proc/meminfo, e.g.:
//		MemTotal:        7628544 kB
//		MemFree:         4038912 kB
//		Buffers:          289732 kB
//		Cached:          1851184 kB
//	2. Split required fields.
//	3. Convert all values to MBytes.
//
// See `man proc` for details.
func (m *Mem) Extract() (results []float64) {
	file, err := os.Open("/proc/meminfo")
	if err != nil {
		log.Println(err)
		return make([]float64, len(m.Columns))
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	for i := 0; i < len(m.Columns); i++ {
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Println(err)
			return make([]float64, len(m.Columns))
		}

		value, err := strconv.ParseFloat(strings.Fields(line)[1], 32)
		value = math.Floor(value / 1024) // KB -> MB
		if err != nil {
			log.Println(err)
			return make([]float64, len(m.Columns))
		}
		results = append(results, value)
	}
	for i := 1; i < len(results); i++ {
		results[0] -= results[i]
	}
	return
}
