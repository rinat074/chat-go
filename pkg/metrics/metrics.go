package metrics

import (
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Metrics содержит все метрики для сервиса.
type Metrics struct {
	registry          *prometheus.Registry
	requestCounter    *prometheus.CounterVec
	requestDuration   *prometheus.HistogramVec
	databaseDuration  *prometheus.HistogramVec
	cacheHits         prometheus.Counter
	cacheMisses       prometheus.Counter
	goroutinesGauge   prometheus.Gauge
	memoryAllocGauge  prometheus.Gauge
	activeConnections prometheus.Gauge
}

// NewMetrics создает новый экземпляр метрик.
func NewMetrics(serviceName string) *Metrics {
	registry := prometheus.NewRegistry()

	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_requests_total", serviceName),
			Help: "Total number of requests",
		},
		[]string{"method", "endpoint", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    fmt.Sprintf("%s_request_duration_seconds", serviceName),
			Help:    "Request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)

	databaseDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    fmt.Sprintf("%s_database_duration_seconds", serviceName),
			Help:    "Database operation duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"operation"},
	)

	cacheHits := prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_cache_hits_total", serviceName),
			Help: "Total number of cache hits",
		},
	)

	cacheMisses := prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_cache_misses_total", serviceName),
			Help: "Total number of cache misses",
		},
	)

	goroutinesGauge := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: fmt.Sprintf("%s_goroutines", serviceName),
			Help: "Current number of goroutines",
		},
	)

	memoryAllocGauge := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: fmt.Sprintf("%s_memory_alloc_bytes", serviceName),
			Help: "Current memory allocations in bytes",
		},
	)

	activeConnections := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: fmt.Sprintf("%s_active_connections", serviceName),
			Help: "Current number of active connections",
		},
	)

	registry.MustRegister(
		requestCounter,
		requestDuration,
		databaseDuration,
		cacheHits,
		cacheMisses,
		goroutinesGauge,
		memoryAllocGauge,
		activeConnections,
	)

	return &Metrics{
		registry:          registry,
		requestCounter:    requestCounter,
		requestDuration:   requestDuration,
		databaseDuration:  databaseDuration,
		cacheHits:         cacheHits,
		cacheMisses:       cacheMisses,
		goroutinesGauge:   goroutinesGauge,
		memoryAllocGauge:  memoryAllocGauge,
		activeConnections: activeConnections,
	}
}

// ServeHTTP предоставляет HTTP обработчик для метрик Prometheus.
func (m *Metrics) ServeHTTP(port int) error {
	http.Handle("/metrics", promhttp.HandlerFor(m.registry, promhttp.HandlerOpts{}))
	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

// ObserveRequest наблюдает за HTTP запросом.
func (m *Metrics) ObserveRequest(method, endpoint, status string, duration time.Duration) {
	m.requestCounter.WithLabelValues(method, endpoint, status).Inc()
	m.requestDuration.WithLabelValues(method, endpoint).Observe(duration.Seconds())
}

// ObserveDatabaseOperation наблюдает за операцией с базой данных.
func (m *Metrics) ObserveDatabaseOperation(operation string, duration time.Duration) {
	m.databaseDuration.WithLabelValues(operation).Observe(duration.Seconds())
}

// IncrementCacheHit увеличивает счетчик попаданий в кэш.
func (m *Metrics) IncrementCacheHit() {
	m.cacheHits.Inc()
}

// IncrementCacheMiss увеличивает счетчик промахов в кэш.
func (m *Metrics) IncrementCacheMiss() {
	m.cacheMisses.Inc()
}

// SetGoroutines устанавливает количество горутин.
func (m *Metrics) SetGoroutines(count int) {
	m.goroutinesGauge.Set(float64(count))
}

// SetMemoryAlloc устанавливает количество выделенной памяти.
func (m *Metrics) SetMemoryAlloc(bytes int64) {
	m.memoryAllocGauge.Set(float64(bytes))
}

// SetActiveConnections устанавливает количество активных соединений.
func (m *Metrics) SetActiveConnections(count int) {
	m.activeConnections.Set(float64(count))
}
