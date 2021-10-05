package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

func main() {
	conf, err := parseArguments(os.Args[1:])
	if err != nil {
		log.Fatal("failed to parse cli arguments:", err)
	}
	c := make(chan metric)
	go loadAvg(c, conf.pollingPeriod)
	go cpuPercent(c, conf.pollingPeriod)
	go memoryUsage(c, conf.pollingPeriod)
	go diskUsage(c, conf.pollingPeriod, conf.partition)
	go networkUsage(c, conf.pollingPeriod, conf.networkInterface)
	go printMetrics(c, conf.pollingPeriod)
	//Todo: allow graceful shutdown
	select {}
}

type configuration struct {
	pollingPeriod    time.Duration
	networkInterface string
	partition        string
}

var defaultConfig = configuration{
	pollingPeriod:    15 * time.Second,
	partition:        "/",
	networkInterface: "eth0",
}

func parseArguments(args []string) (configuration, error) {
	conf := defaultConfig
	if len(args) == 0 {
		return conf, nil
	}
	if len(args)%2 != 0 {
		return conf, fmt.Errorf("invalid number of arguments")
	}
	for i := 0; i < len(args); i += 2 {
		switch args[i] {
		case "-i":
			d, err := strconv.Atoi(args[i+1])
			if err != nil {
				return conf, err
			}
			conf.pollingPeriod = time.Duration(d) * time.Second
		case "-n":
			conf.networkInterface = args[i+1]
		case "-p":
			conf.partition = args[i+1]
		default:
			return conf, fmt.Errorf("unknown option %s", args[0])
		}
	}
	return conf, nil
}

type metric struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type metricsOutput struct {
	//Timestamp is the unix epoch
	Timestamp int64    `json:"timestamp"`
	Metrics   []metric `json:"metrics"`
}

//printMetrics periodically prints the latest metrics in stdout
func printMetrics(metrics chan metric, outputPeriod time.Duration) {
	ticker := time.NewTicker(outputPeriod)
	curentMetrics := map[string]metric{}
	for {
		select {
		case <-ticker.C:
			out := metricsOutput{
				Timestamp: time.Now().UTC().Unix(),
				Metrics:   toSlice(curentMetrics),
			}
			b, err := json.Marshal(out)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(string(b))
		case m := <-metrics:
			//Note: metrics' names must be unique
			curentMetrics[m.Name] = m
		}
	}
}

//toSlice converts a map of metrics in a slice. The slice order is
//random.
func toSlice(metrics map[string]metric) []metric {
	r := []metric{}
	for _, m := range metrics {
		r = append(r, m)
	}
	return r
}
