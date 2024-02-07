package main

import (
	"test/actuator"
	"test/sensor"
)

func main() {
	f := make(chan struct{})
	go sensor.ServeSensorServer()
	go actuator.ServeActuatorServer()
	<-f
}
