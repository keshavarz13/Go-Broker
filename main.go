package main

import (
	apiImp "therealbroker/api/Imp"
	"therealbroker/pkg/prometheus"
)

func main() {
	go prometheus.AddPrometheus()
	apiImp.RunGrpcServer()
}
