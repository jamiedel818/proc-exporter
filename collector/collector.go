package collector

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type Collector interface {
	ParseProcFile(procFileName string) error
	GetMetricData() map[string]uint64
	// Probably a method that can output prom style metrics
}

// ---------------------- //
// MemInfo Implementation //
// ---------------------- //

// Represents metrics pulled from /proc/meminfo
type MemInfo struct {
	data map[string]uint64
}

// ParseProcFile opens and closes proc file. Sets objects unexported `data“ field with the output of parseMemInfo().
func (m *MemInfo) ParseProcFile(procFileName string) error {
	fi, err := os.Open(procFileName)
	if err != nil {
		return fmt.Errorf("could not open meminfo proc file %q. %w", procFileName, err)
	}

	defer fi.Close()

	memInfo, err := parseMemInfo(fi)
	if err != nil {
		return fmt.Errorf("could not parse meminfo. %w", err)
	}

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
