package main

import (
	"log"
	apiImp "therealbroker/api/Imp"
	"therealbroker/pkg/prometheus"
)

func main() {
	log.Println("Broker started!")
	go prometheus.AddPrometheus()
	apiImp.RunGrpcServer()
}
