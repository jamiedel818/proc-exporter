package collector

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOutputHelp(t *testing.T) {
	m := promMetric{
		name:        "test",
		metricType:  metricGuage,
		description: "test",
	}

	assert.Equal(t, "# HELP test test\n", m.outputHelp())
}

func TestOutputType(t *testing.T) {
	m := promMetric{
		name:        "test",
		metricType:  metricGuage,
		description: "test",
	}

	assert.Equal(t, "# TYPE test gauge\n", m.outputType())
}

func TestIsValidMetric(t *testing.T) {
	tests := []struct {
		name   string
		metric promMetric
		want   bool
	}{
		{
			name: "valid gauge metric",
			metric: promMetric{
				name:       "memory_usage",
				metricType: metricGuage,
			},
			want: true,
		},
		{
			name: "valid counter metric",
			metric: promMetric{
				name:       "http_requests_total",
				metricType: metricCounter,
			},
			want: true,
		},
		{
			name: "valid histogram metric",
			metric: promMetric{
				name:       "request_duration_seconds",
				metricType: metricHistorgram,
			},
			want: true,
		},
		{
			name: "invalid metric type",
			metric: promMetric{
				name:       "memory_usage",
				metricType: "invalid_type",
			},
			want: false,
		},
	}

	for _, c := range tests {
		t.Run(c.name, func(t *testing.T) {
			assert.Equal(t, c.want, c.metric.isValidateMetric())
		})
	}
}
