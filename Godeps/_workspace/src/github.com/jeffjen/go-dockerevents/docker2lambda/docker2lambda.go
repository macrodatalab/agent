package main

import (
	d2k "github.com/yihungjen/go-dockerevents"
	"log"
)

func main() {
	log.Println("docker2lambda monitor begin...")

	// Begin monitoring Docker Event
	lambdaSink := d2k.EventLoop(nil, 100)

	// Start driving docker events to AWS lambda function
	LambdaLoop(lambdaSink, d2k.DefaultCallBack)
}
