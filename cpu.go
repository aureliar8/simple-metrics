package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	"log"
)

//getLoadAverage read /proc/loadavg and return three metrics for the
//load average in the last 1, 5 and 15 minutes.
func getLoadAverage() ([]metric, error) {
	b, err := ioutil.ReadFile("/proc/loadavg")
	if err != nil {
		return nil, err
	}
	fields := strings.Fields(string(b))
	if len(fields) < 3 {
		return nil, fmt.Errorf("unexpected /proc/loadavg format: got %d fields", len(fields))
	}
	return []metric{{
		Name:  "loadavg-1",
		Value: fields[0],
	}, {
		Name:  "loadavg-5",
		Value: fields[1],
	}, {
		Name:  "loadavg-15",
		Value: fields[2],
	},
	}, nil
}

//getCPUPercentage get the usage percentage for each cpu (and in
//average). This is calculated since the boot. Todo: Only emit usage for a recent time
func getCPUPercentage() ([]metric, error) {
	b, err := ioutil.ReadFile("/proc/stat")
	if err != nil {
		return nil, err
	}
	metrics := []metric{}
	lines := strings.Split(string(b), "\n")
	for _, line := range lines {
		if !strings.HasPrefix(line, "cpu") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 5 {
			return nil, fmt.Errorf("invalid /proc/stat format")
		}
		user, err := strconv.Atoi(fields[1])
		if err != nil {
			return nil, err
		}
		nice, err := strconv.Atoi(fields[2])
		if err != nil {
			return nil, err
		}
		system, err := strconv.Atoi(fields[3])
		if err != nil {
			return nil, err
		}
		idle, err := strconv.Atoi(fields[4])
		if err != nil {
			return nil, err
		}
		usage := float64(user+nice+system) / float64(user+nice+system+idle)
		metrics = append(metrics, metric{
			Name:  fields[0],
			Value: fmt.Sprintf("%f", usage*100),
		})
	}
	return metrics, nil
}

//cpuPercent periodically pools the cpu usage and writes the
//metrics to the channel.
func cpuPercent(out chan metric, period time.Duration) {
	ticker := time.NewTicker(period)
	for {
		metrics, err := getCPUPercentage()
		if err != nil {
			log.Println(fmt.Errorf("cpu metrics: %w", err))
		}
		for _, m := range metrics {
			out <- m
		}
		<-ticker.C
	}
}

//loadAvg periodically pools the load average and writes the metrics to
//the channel.
func loadAvg(out chan metric, period time.Duration) {
	ticker := time.NewTicker(period)
	for {
		metrics, err := getLoadAverage()
		if err != nil {
			log.Println(fmt.Errorf("cpu loadavg: %w", err))
		}
		for _, m := range metrics {
			out <- m
		}
		<-ticker.C
	}
}
