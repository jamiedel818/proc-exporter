package collector

import (
	"fmt"
	"regexp"
	"slices"
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

// Collector represents the common operations to retrieve data from a linux proc file and output the data
// in prometheus style metrics. How this is done is up to the individual implementation
type Collector interface {
	ParseProcFile() error
	OutputPromMetrics() string
}

// promMetric aids in the production of prometheus style metrics.
type promMetric struct {
	name        string
	metricType  string
	description string
}

/*
isValidateMetric checks whether the given Prometheus metric configuration is valid.
The main checks are whether the provided metric is a valid Prometheus metric type
and if it matches the documented naming format provided my Prometheus - https://prometheus.io/docs/concepts/data_model/
*/
func (m promMetric) isValidateMetric() bool {
	var ok bool

	if ok = slices.Contains(validMetricTypes, m.metricType); !ok {
		return ok
	}

	if ok, _ = regexp.MatchString(promMetricNameRegex, m.name); !ok {
		return ok
	}

	return ok
}

func (m promMetric) outputHelp() string {
	return fmt.Sprintf("# HELP %s %s\n", m.name, m.description)
}

func (m promMetric) outputType() string {
	return fmt.Sprintf("# TYPE %s %s\n", m.name, m.metricType)
}
