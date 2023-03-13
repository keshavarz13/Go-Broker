package main

import (
	apiImp "therealbroker/api/Imp"
	modules "therealbroker/modules"
)

func main() {
	go modules.AddPrometheus()
	apiImp.RunGrpcServer()
}
