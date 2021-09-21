package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"time"
)

//getMemoryUsage return a metrics with the current percentage of
//memory used by the machine
func getMemoryUsage() (metric, error) {
	m := metric{Name: "memory usage"}
	b, err := ioutil.ReadFile("/proc/meminfo")
	if err != nil {
		return m, err
	}
	lines := strings.Split(string(b), "\n")
	available := -1
	total := -1
	for _, line := range lines {
		if strings.HasPrefix(line, "MemAvailable") {
			fields := strings.Fields(line)
			value, err := strconv.Atoi(fields[1])
			if err != nil {
				return m, err
			}
			factor, err := unitFactor(fields[2])
			if err != nil {
				return m, nil
			}
			available = value * factor
		}
		if strings.HasPrefix(line, "MemTotal") {
			fields := strings.Fields(line)
			value, err := strconv.Atoi(fields[1])
			if err != nil {
				return m, err
			}
			factor, err := unitFactor(fields[2])
			if err != nil {
				return m, nil
			}
			total = value * factor
		}
	}
	usage := float64(1) - float64(available)/float64(total)
	m.Value = fmt.Sprintf("%f", usage*100)
	return m, nil
}

//memoryUsage periodical pools the memory usage and write the metric to the channel
func memoryUsage(out chan metric, period time.Duration) {
	ticker := time.NewTicker(period)
	for {
		m, err := getMemoryUsage()
		if err != nil {
			log.Println(fmt.Errorf("memoru usage: %w", err))
		}
		out <- m
		<-ticker.C
	}
}

//unitFactor return the unitFactor of the given unit.
func unitFactor(unit string) (int, error) {
	switch unit {
	case "B":
		return 1, nil
	case "kB":
		return 1000, nil
	case "MB":
		return 1_000_000, nil
	default:
		return 0, fmt.Errorf("unknow unit %s", unit)
	}
}
