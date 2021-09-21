package main

import (
	"fmt"
	"log"
	"syscall"
	"time"
)

//getPartitionUsage return a metrics with the disk usage of the given partition
func getPartitionUsage(path string) (metric, error) {
	m := metric{Name: path + " usage"}
	fs := syscall.Statfs_t{}
	err := syscall.Statfs(path, &fs)
	if err != nil {
		return m, err
	}
	all := fs.Blocks * uint64(fs.Bsize)
	available := fs.Bavail * uint64(fs.Bsize)
	usage := float64(all-available) / float64(all)
	m.Value = fmt.Sprintf("%f", usage*100)
	return m, nil
}

//diskUsage periodical pools the disk usage and writes the metric in the channel
func diskUsage(out chan metric, period time.Duration, partition string) {
	ticker := time.NewTicker(period)
	for {
		m, err := getPartitionUsage(partition)
		if err != nil {
			log.Println(fmt.Errorf("disk %s: %w", partition, err))
		}
		out <- m
		<-ticker.C
	}
}
