package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// HTTP metrics
	httpRequestsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "myapp_http_requests_total",
		Help: "Total number of HTTP requests",
	}, []string{"service", "handler", "method", "code"})

	httpRequestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "myapp_http_request_duration_seconds",
		Help:    "Duration of HTTP requests",
		Buckets: []float64{0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10},
	}, []string{"service", "handler", "method"})

	httpRequestsInFlight = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "myapp_http_requests_in_flight",
		Help: "Current number of HTTP requests being served",
	})

	// Business metrics
	ordersProcessed = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "myapp_orders_processed_total",
		Help: "Total number of processed orders",
	}, []string{"service", "status"})

	// Custom system metrics (using custom names to avoid conflicts)
	appGoroutines = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "myapp_custom_goroutines",
		Help: "Custom goroutine count measurement",
	})

	appMemoryAlloc = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "myapp_custom_memory_alloc_bytes",
		Help: "Custom memory allocation measurement",
	})
)

func recordSystemMetrics(ctx context.Context) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			appGoroutines.Set(float64(runtime.NumGoroutine()))
			var m runtime.MemStats
			runtime.ReadMemStats(&m)
			appMemoryAlloc.Set(float64(m.Alloc))
		case <-ctx.Done():
			return
		}
	}
}

func instrumentHandler(serviceName, handlerName string, handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		httpRequestsInFlight.Inc()
		defer httpRequestsInFlight.Dec()

		rww := &responseWriterWrapper{ResponseWriter: w}

		handler(rww, r)

		duration := time.Since(start).Seconds()
		statusCode := fmt.Sprintf("%d", rww.statusCode)

		httpRequestsTotal.WithLabelValues(
			serviceName,
			handlerName,
			r.Method,
			statusCode,
		).Inc()

		httpRequestDuration.WithLabelValues(
			serviceName,
			handlerName,
			r.Method,
		).Observe(duration)
	}
}

type responseWriterWrapper struct {
	http.ResponseWriter
	statusCode int
}

func (rww *responseWriterWrapper) WriteHeader(code int) {
	rww.statusCode = code
	rww.ResponseWriter.WriteHeader(code)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func orderHandler(w http.ResponseWriter, r *http.Request) {
	////time.Sleep(time.Duration(50+rand.Intn(100)) * time.Millisecond)

	if rand.Float32() > 0.1 {
		ordersProcessed.WithLabelValues("order-service", "success").Inc()
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "processed"}`))
	} else {
		ordersProcessed.WithLabelValues("order-service", "failed").Inc()
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "processing failed"}`))
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start system metrics collection
	go recordSystemMetrics(ctx)

	// Set up router
	mux := http.NewServeMux()

	// Instrumented handlers
	mux.HandleFunc("/health", instrumentHandler("order-service", "health", healthHandler))
	mux.HandleFunc("/order", instrumentHandler("order-service", "order", orderHandler))

	// Prometheus metrics endpoint
	mux.Handle("/metrics", promhttp.Handler())

	fmt.Println("Server Start in 0.0.0.0:8080")
	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	// Shutdown gracefully
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("Graceful shutdown failed: %v", err)
	}
}
