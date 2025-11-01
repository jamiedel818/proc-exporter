package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

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
	Lock    sync.Mutex
	Metrics []collector.Collector
}

func (h *AppHandler) CollectMetrics() {
	ticker := time.NewTicker(5 * time.Second)

	defer ticker.Stop()
	for {
		<-ticker.C

		fmt.Println("Taking the lock to write data")
		h.Lock.Lock()

		for _, c := range h.Metrics {
			err := c.ParseProcFile()
			if err != nil {
				fmt.Println(err)
			}
		}
		fmt.Println("Lock released after writing data")
		h.Lock.Unlock()
	}
}

func (h *AppHandler) GetMetrcs() []map[string]uint64 {
	fmt.Println("Aquiring the lock to read metrics")
	h.Lock.Lock()
	defer h.Lock.Unlock()

	result := []map[string]uint64{}
	for _, c := range h.Metrics {
		result = append(result, c.GetMetricData())
	}
	fmt.Println("Releasing the lock after reading metrics")
	return result
}

func main() {
	handler := &AppHandler{
		Lock: sync.Mutex{},
		Metrics: []collector.Collector{
			&collector.MemInfo{ProcFileName: "./fixtures/meminfo_full"},
		},
	}

	go handler.CollectMetrics()

	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "%v", handler.GetMetrcs())
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
