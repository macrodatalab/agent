package main

import (
	d2k "github.com/yihungjen/go-dockerevents"
	"log"
)

func main() {
	log.Print("docker2kinesis monitor begin...")

	// Begin monitoring Docker Event
	ksisSink := d2k.EventLoop(nil, 100)

	// Start driving docker events to AWS kinesis
	KinesisLoop(ksisSink, d2k.DefaultCallBack)
}
