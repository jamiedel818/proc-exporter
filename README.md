# proc-exporter

A lightweight option to export metrics from proc files in Prometheus format

## Usage
Start the exporter with default settings:

```
proc-exporter -d <path to proc directory> -p <server port> -i <scrape interval>
```

Example:
```
proc-exporter -d /proc -p 9100 -i 15
```

## Adding Metrics
To add more proc files or metrics, simply implement the `Collector` interface, and add it to the `AppHandler`
