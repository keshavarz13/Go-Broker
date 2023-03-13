package modules

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func AddPrometheus() {
	http.Handle("/metrics", promhttp.Handler())
	err := http.ListenAndServe(":9000", nil)
	log.Fatal(err)
}
