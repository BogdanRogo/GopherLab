package utils

import (
	"log"
	"time"

	"go.opencensus.io/exporter/prometheus"
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
)

var (
	// MLatencyMs The latency in milliseconds
	MLatencyMs = stats.Float64("urlshortener/latency", "The latency in milliseconds per short request", "ms")
	// MErrors Encounters the number of non EOF(end-of-file) errors.
	MErrors = stats.Int64("urlshortener/errors", "The number of errors encountered", "1")
	// HTTPMethod ...
	HTTPMethod, _ = tag.NewKey("method")
	// HTTPHandler ...
	HTTPHandler, _ = tag.NewKey("service")
	// HTTPStatus ...
	HTTPStatus, _ = tag.NewKey("status")
)

var (
	LatencyView = &view.View{
		Name:        "urlshortener/latency",
		Measure:     MLatencyMs,
		Description: "The distribution of the latencies",

		// Latency in buckets:
		// [>=5ms, >=10ms ... >=4s, >=6s]
		Aggregation: view.Distribution(5, 10, 15, 20, 25, 50, 75, 100, 250, 500, 750, 1000, 5000),
		TagKeys:     []tag.Key{HTTPMethod, HTTPHandler, HTTPStatus}}

	ErrorCountView = &view.View{
		Name:        "urlshortener/errors",
		Measure:     MErrors,
		Description: "The number of errors encountered",
		Aggregation: view.Count(),
	}
)

// OCPrometheusExporter exports an OpenCensus prometheus Exporter
func OCPrometheusExporter() (Exporter *prometheus.Exporter) {
	exporter, err := prometheus.NewExporter(prometheus.Options{Namespace: "reliability"})
	if err != nil {
		log.Fatal(err)
	}
	view.RegisterExporter(exporter)

	if err := view.Register(LatencyView, ErrorCountView); err != nil {
		log.Fatalf("Failed to register views: %v", err)
	}
	view.SetReportingPeriod(1 * time.Second)

	return exporter
}
