package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"
)

//getNeworkStats returns two metrics which are the number of bytes
//send and received on the given interface.
func getNeworkStats(netInterface string) ([]metric, error) {
	b, err := ioutil.ReadFile("/proc/net/dev")
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(b), "\n")
	for _, line := range lines {

		fields := strings.Fields(line)
		if len(fields) == 0 {
			//Do not count empty line as an error
			if line == "" {
				continue
			}
			return nil, fmt.Errorf("Invalid /proc/net/dev format")
		}
		if fields[0] == netInterface+":" {
			if len(fields) < 10 {
				return nil, fmt.Errorf("Invalid /proc/net/dev format")
			}
			return []metric{
				{Name: netInterface + " tx", Value: fields[1]},
				{Name: netInterface + " rx", Value: fields[9]},
			}, nil
		}
	}
	return nil, fmt.Errorf("network interface %s not found", netInterface)
}

//networkUsage periodical pools network statistics about the given
//interface and write the metrics to the channel.
func networkUsage(out chan metric, period time.Duration, netInterface string) {
	ticker := time.NewTicker(period)
	for {
		metrics, err := getNeworkStats(netInterface)
		if err != nil {
			log.Println(fmt.Errorf("%s network usage: %w ", netInterface, err))
		}
		for _, m := range metrics {
			out <- m
		}
		<-ticker.C
	}
}
