package prometheus

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	ActiveSubscribers = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "broker_active_subscribers",
		Help: "number of active subscribers in broker",
	})
	MethodCalls = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "method_count",
		Help: "number of method calls in broker",
	}, []string{"method"})

	MethodError = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "method_error_count",
		Help: "counter error of each method",
	}, []string{"method"})
)

func AddPrometheus() {
	http.Handle("/metrics", promhttp.Handler())
	err := http.ListenAndServe(":9000", nil)
	log.Fatal(err)
}
