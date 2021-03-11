package main

import (
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"strconv"
	"time"
)

var (
	prom_version = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "version",
		Help: "Version information about this binary",
		ConstLabels: map[string]string{
			"version": "v0.0.2",
		},
	})

	prom_httpRequestTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Count of all http requests",
	}, []string{"method", "code"})

	prom_httpRequestDurationSeconds = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration distribution",
			Buckets: []float64{1, 2, 5, 10, 20, 60},
		}, []string{"method"})
)

func getRoute(res http.ResponseWriter, req *http.Request) {
	start := time.Now()
	duration := time.Since(start)
	res.Write([]byte("home route!"))
	res.WriteHeader(http.StatusOK)
	prom_httpRequestDurationSeconds.With(prometheus.Labels{"method": "GET"}).Observe(duration.Seconds())
	prom_httpRequestTotal.With(prometheus.Labels{"method": "GET", "code": strconv.Itoa(http.StatusOK)}).Inc()
}

func postRoute(res http.ResponseWriter, req *http.Request) {
	start := time.Now()
	duration := time.Since(start)
	res.Write([]byte("ok route!"))
	res.WriteHeader(http.StatusCreated)
	prom_httpRequestDurationSeconds.With(prometheus.Labels{"method": "POST"}).Observe(duration.Seconds())
	prom_httpRequestTotal.With(prometheus.Labels{"method": "POST", "code": strconv.Itoa(http.StatusNoContent)}).Inc()
}

func main() {
	router := mux.NewRouter().StrictSlash(true)

	prom := prometheus.NewRegistry()
	prom.MustRegister(prom_version)
	prom.MustRegister(prom_httpRequestTotal)
	prom.MustRegister(prom_httpRequestDurationSeconds)

	router.HandleFunc("/", getRoute).Methods("GET")
	router.HandleFunc("/index", postRoute).Methods("POST")
	router.Handle("/metrics", promhttp.HandlerFor(prom, promhttp.HandlerOpts{}))
	log.Println(http.ListenAndServe(":8080", router))
}
