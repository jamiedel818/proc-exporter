package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/jamidel818/mini-node-exporter/collector"
)

// AppHandler is the high level struct responsible for storing shared data and orchestrating the exporter.
type AppHandler struct {
	Lock           sync.Mutex
	Metrics        []collector.Collector
	ScrapeInterval int
}

// CollectMetrics continuously instructs the configured collector.Collectors to get the latest metric data.
// Data is stored in each collectors unexported data field.
// TODO make the ticker configurable
func (h *AppHandler) CollectMetrics() {
	ticker := time.NewTicker(time.Duration(h.ScrapeInterval) * time.Second)

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
	procDir := flag.String("d", "/proc", "path to the systems proc directory")
	port := flag.Int("p", 8080, "port to run the http server")
	interval := flag.Int("i", 10, "interval to scrape proc files")

	flag.Parse()
	fmt.Println("Starting proc-exporter")
	fmt.Printf("PROC DIR: %s\n", *procDir)
	fmt.Printf("PORT: %d\n", *port)
	fmt.Printf("INTERVAL: %d\n", *interval)

	handler := &AppHandler{
		Lock:           sync.Mutex{},
		ScrapeInterval: *interval,
		Metrics: []collector.Collector{
			&collector.MemInfo{ProcFileName: fmt.Sprintf("%s/meminfo", *procDir)},
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

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
}
