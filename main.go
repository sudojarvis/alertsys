package main

import (
	"net/http"
	"os"

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
)

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		timer := prometheus.NewTimer(httpDuration.WithLabelValues(r.URL.Path))
		defer timer.ObserveDuration()

		httpRequestsTotal.WithLabelValues(r.URL.Path).Inc()
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Welcome to the dummy service!"))
	})

	// http.HandleFunc("/alert", handleAlert)

	http.Handle("/metrics", promhttp.Handler())

	log.Info("Starting HTTP server on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Error starting HTTP server: %v", err)
	}
}
