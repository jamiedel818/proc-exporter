package collector

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

const (
	metricGuage      = "gauge"
	metricCounter    = "counter"
	metricHistorgram = "histogram"

	promMetricNameRegex = "[a-zA-Z_][a-zA-Z0-9_]*"
)

var validMetricTypes = []string{
	metricGuage,
	metricCounter,
	metricHistorgram,
}

var memInfoToProm = map[string]promMetric{
	"memavailable": {
		name:        "memory_available_bytes",
		mType:       metricGuage,
		description: "An estimate of how much memory is available for new applications without triggering the system to swap.",
	},
	"memtotal": {
		name:        "memory_total_bytes",
		mType:       metricGuage,
		description: "Total usable physical memory on the device.",
	},
}

type Collector interface {
	ParseProcFile() error
	GetMetricData() map[string]uint64
	OutputPromMetrics() string
}

type promMetric struct {
	name        string
	mType       string
	description string
}

// ---------------------- //
// MemInfo Implementation //
// ---------------------- //

// Represents metrics pulled from /proc/meminfo
type MemInfo struct {
	ProcFileName string
	data         map[string]uint64
}

// ParseProcFile opens and closes proc file. Sets objects unexported `data“ field with the output of parseMemInfo().
func (m *MemInfo) ParseProcFile() error {
	fi, err := os.Open(m.ProcFileName)
	if err != nil {
		return fmt.Errorf("could not open meminfo proc file %q. %w", m.ProcFileName, err)
	}

	defer fi.Close()

	memInfo, err := parseMemInfo(fi)
	if err != nil {
		return fmt.Errorf("could not parse meminfo. %w", err)
	}

	// convert to bytes
	for k, v := range memInfo {
		memInfo[k] = v * 1024
	}

	m.data = memInfo

	return nil
}

// GetMetricData retrieves objects unexported `data` field
func (m *MemInfo) GetMetricData() map[string]uint64 {
	return m.data
}

// parseMemInfo expects io.Reader in the form of a /proc/meminfo file.
// Parses out and normalizes all metrics and their values. Values are returned in their original form, kB (1024 bytes)
func parseMemInfo(r io.Reader) (map[string]uint64, error) {
	sc := bufio.NewScanner(r)

	d := make(map[string]uint64)
	for sc.Scan() {
		// example line: "MemTotal:       16397740 kB"
		l := sc.Text()
		k, v, ok := strings.Cut(l, ":")
		if !ok {
			// seperator was not found. skip the line
			continue
		}

		v = strings.TrimSuffix(v, "kB")
		v = strings.TrimSpace(v)

		n, err := strconv.ParseUint(v, 10, 64)

		if err != nil {
			return map[string]uint64{}, fmt.Errorf("could not convert %q to uint64 for metric %q. %w", v, k, err)
		}

		d[strings.TrimSpace(strings.ToLower(k))] = n
	}

	if err := sc.Err(); err != nil {
		return map[string]uint64{}, fmt.Errorf("bufio scan error occured. %w", err)
	}

	return d, nil
}

// memory_available_bytes
// memory_total_bytes

func (m *MemInfo) OutputPromMetrics() string {
	// we capture the whole meminfo file. Here we can output the metrics we care about
	promMetrics := ""
	for k, metric := range memInfoToProm {
		if !metric.isValidateMetric() {
			fmt.Printf("invalid prometheus metric configuration provided %#v", m)
			continue
		}

		promMetrics += outputHelp(metric.name, metric.description) + outputType(metric.name, metric.mType)
		promMetrics += fmt.Sprintf("%d\n", m.data[k])

	}
	return promMetrics
}

func outputHelp(mName string, mDesc string) string {
	return fmt.Sprintf("# HELP %s %s\n", mName, mDesc)
}

func outputType(mName string, mType string) string {
	return fmt.Sprintf("# TYPE %s %s\n", mName, mType)
}

/*
isValidateMetric checks whether the given Prometheus metric configuration is valid.
The main checks are whether the provided metric is a valid Prometheus metric type
and if it matches the documented naming format provided my Prometheus - https://prometheus.io/docs/concepts/data_model/
*/
func (m promMetric) isValidateMetric() bool {
	var ok bool

	if ok = slices.Contains(validMetricTypes, m.mType); !ok {
		return ok
	}

	if ok, _ = regexp.MatchString(promMetricNameRegex, m.name); !ok {
		return ok
	}

	return ok
}
