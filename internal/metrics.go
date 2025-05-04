package internal

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	extractionsTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "etl_extractions_total",
		Help: "Total number of successful data extractions",
	})

	transformationsTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "etl_transformations_total",
		Help: "Total number of successful data transformations",
	})

	loadsTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "etl_loads_total",
		Help: "Total number of successful data loads",
	})

	failuresTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "etl_failures_total",
		Help: "Total number of ETL failures",
	})
)

func InitMetrics() error {
	// Register metrics
	prometheus.MustRegister(extractionsTotal)
	prometheus.MustRegister(transformationsTotal)
	prometheus.MustRegister(loadsTotal)
	prometheus.MustRegister(failuresTotal)

	// Start HTTP server
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	go func() {
		if err := http.ListenAndServe(":2112", nil); err != nil {
			fmt.Printf("Failed to start metrics server: %v\n", err)
		}
	}()

	return nil
}

func IncrementExtractions() {
	extractionsTotal.Inc()
}

func IncrementTransformations() {
	transformationsTotal.Inc()
}

func IncrementLoads() {
	loadsTotal.Inc()
}

func IncrementFailures() {
	failuresTotal.Inc()
}
