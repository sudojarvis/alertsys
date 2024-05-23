package main

import (
	"net/http"
	"os"
	"sync"
	"errors"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

// type Alert struct {
// 	Status string `json:"status"`
// 	Labels struct {
// 		Alertname string `json:"alertname"`
// 		Instance  string `json:"instance"`
// 		Job       string `json:"job"`
// 	} `json:"labels"`
// }

// func handleAlert(w http.ResponseWriter, r *http.Request) {
// 	var alert Alert
// 	err := json.NewDecoder(r.Body).Decode(&alert)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	if alert.Status == "firing" {
// 		if alert.Labels.Job == "go-service" {
// 			// Restart the Go service
// 			cmd := exec.Command("sh", "-c", "docker restart go-webserver")
// 			err := cmd.Run()
// 			if err != nil {
// 				log.Printf("Failed to restart Go service: %v", err)
// 				http.Error(w, "Failed to restart Go service", http.StatusInternalServerError)
// 				return
// 			}
// 			log.Println("Go service restarted")
// 		}
// 	}

// 	w.WriteHeader(http.StatusOK)
// }

var (
    scannerHealthy bool = true // Health status of the scanner
    mu             sync.RWMutex // Mutex to protect the health status variable
)

// Define Prometheus metrics
var (
    httpRequestsTotal = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total number of HTTP requests processed",
        },
        []string{"path"},
    )
    httpDuration = promauto.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "http_request_duration_seconds",
            Help:    "Duration of HTTP requests",
            Buckets: prometheus.DefBuckets,
        },
        []string{"path"},
    )
    scanFailures = prometheus.NewCounter(
        prometheus.CounterOpts{
            Name: "scanner_failures_total",
            Help: "Total number of scanner failures.",
        },
    )
)

// init function to register metrics
func init() {
    prometheus.MustRegister(scanFailures) // Register the scanner failures metric explicitly
}

func main() {
    // Configure logger
    log.SetFormatter(&log.JSONFormatter{})
    log.SetOutput(os.Stdout)

    // Root endpoint handler
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        timer := prometheus.NewTimer(httpDuration.WithLabelValues(r.URL.Path))
        defer timer.ObserveDuration()

        httpRequestsTotal.WithLabelValues(r.URL.Path).Inc()
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("Welcome to the dummy service!"))
    })

    // Prometheus metrics endpoint
    http.Handle("/metrics", promhttp.Handler())

    // Health check endpoint
    http.HandleFunc("/health", healthCheckHandler)

    // Start scanner in a goroutine
    go func() {
        for {
            if err := runScanner(); err != nil {
                scanFailures.Inc()
                log.Printf("Scanner failed: %v", err)
                setScannerHealthy(false)
            } else {
                setScannerHealthy(true)
            }
            time.Sleep(1 * time.Minute) // Adjust the interval as needed
        }
    }()

    // Start HTTP server
    log.Info("Starting HTTP server on port 8080")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatalf("Error starting HTTP server: %v", err)
    }
}

// Simulate scanner work
func runScanner() error {
    // Simulate scanner work
    // return nil // Uncomment to simulate successful scanner run
    return errors.New("simulated scanner failure") // Simulate a scanner failure
}

// Health check handler
func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
    if isScannerHealthy() {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("OK"))
    } else {
        w.WriteHeader(http.StatusServiceUnavailable)
        w.Write([]byte("Scanner Unhealthy"))
    }
}

// Set scanner health status
func setScannerHealthy(healthy bool) {
    mu.Lock()
    defer mu.Unlock()
    scannerHealthy = healthy
}

// Get scanner health status
func isScannerHealthy() bool {
    mu.RLock()
    defer mu.RUnlock()
    return scannerHealthy
}