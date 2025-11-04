package collector

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// Represents metrics pulled from /proc/meminfo
type MemInfo struct {
	ProcFileName string
	data         map[string]uint64
}

// Maps a metrics pulled from meminfo to metadata for prometheus style metrics.
// Allows this implementation to define which metrics will be exposed and how they will be displayed.
var memInfoToProm = map[string]promMetric{
	"memavailable": {
		name:        "memory_available_bytes",
		metricType:  metricGuage,
		description: "An estimate of how much memory is available for new applications without triggering the system to swap.",
	},
	"memtotal": {
		name:        "memory_total_bytes",
		metricType:  metricGuage,
		description: "Total usable physical memory on the device.",
	},
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

// OutputPromMetrics formats and returns a string of metrics in prometheus format
func (m *MemInfo) OutputPromMetrics() string {
	var result string

	for metric, pMetric := range memInfoToProm {
		if !pMetric.isValidateMetric() {
			fmt.Printf("invalid prometheus metric configuration provided %#v", m)
			continue
		}

		result += pMetric.outputHelp() + pMetric.outputType()
		result += fmt.Sprintf("%d\n", m.data[metric])

	}
	return result
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
