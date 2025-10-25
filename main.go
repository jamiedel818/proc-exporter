package main

import (
	"fmt"

	"github.com/jamidel818/mini-node-exporter/collector"
)

/*
HTTP serving and metrics collection should be their own goroutines. HTTP can be main

Algo:
Start a timer(stored in main object)
everytime the timer hits
	- aquire a lock stored in the main object
	- open files sequentially (mem, cpu, etc. This can be a []ProcFile{})
	- update the data, release the lock

HTTP server
- Try and aquire the lock on the main object to read the data. Generate response on /metrics
  in prom style metrics. Maybe each ProcFile has a method for generating this.
- release the lock

Maybe use basic system logging

*/

type AppHandler struct {
	// timer
	// lock
	Metrics []collector.Collector
}

func main() {
	handler := AppHandler{
		Metrics: []collector.Collector{
			&collector.MemInfo{},
		},
	}

	for _, m := range handler.Metrics {
		fmt.Printf("Data before parsing %v\n", m.GetMetricData())
		fmt.Println("----------")
		err := m.ParseProcFile("./fixtures/meminfo_full")
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("Data after parsing %v\n", m.GetMetricData())
	}
}
