package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/jamidel818/mini-node-exporter/collector"
)

// AppHandler is the high level struct responsible for storing shared data and orchestrating the exporter.
type AppHandler struct {
	Lock    sync.Mutex
	Metrics []collector.Collector
}

// CollectMetrics continuously instructs the configured collector.Collectors to get the latest metric data.
// Data is stored in each collectors unexported data field.
// TODO make the ticker configurable
func (h *AppHandler) CollectMetrics() {
	ticker := time.NewTicker(5 * time.Second)

	defer ticker.Stop()
	for {
		<-ticker.C

		h.Lock.Lock()

		for _, c := range h.Metrics {
			err := c.ParseProcFile()
			if err != nil {
				fmt.Println(err)
			}
		}

		h.Lock.Unlock()
	}
}

// GetMetrcs instructs the configured collector.Collectors to output their stored metrics in prometheus format.
func (h *AppHandler) GetMetrcs() string {
	h.Lock.Lock()
	defer h.Lock.Unlock()

	metrics := ""
	for _, c := range h.Metrics {
		metrics += c.OutputPromMetrics()
	}

	return metrics
}

func main() {
	fmt.Println("starting mini-node-exporter")

	handler := &AppHandler{
		Lock: sync.Mutex{},
		Metrics: []collector.Collector{
			&collector.MemInfo{ProcFileName: "./fixtures/meminfo_full"}, // TODO this needs to be the real procfiles location
		},
	}

	go handler.CollectMetrics()

	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "%v", handler.GetMetrcs())
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "ok")
	})

	// TODO make the port configurable
	log.Fatal(http.ListenAndServe(":8080", nil))
}
