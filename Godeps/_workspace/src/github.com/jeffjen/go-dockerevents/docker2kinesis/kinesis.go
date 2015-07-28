package main

import (
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/kinesis"
	d2k "github.com/yihungjen/go-dockerevents"
	"log"
	"os"
	"time"
)

var (
	KinesisStreamName string
)

func init() {
	KinesisStreamName = os.Getenv("AWS_KINESIS_STREAM_NAME")
}

func KinesisLoop(event <-chan *d2k.Event, callback func(*d2k.Event) interface{}) {
	// main driver for pushing events to Kinesis

	svc := kinesis.New(&aws.Config{Region: "us-west-2"})

	// Wait... for stream to be ready
	for {
		params := &kinesis.DescribeStreamInput{
			StreamName: &KinesisStreamName,
		}
		if resp, err := svc.DescribeStream(params); err != nil {
			if awsErr, ok := err.(awserr.Error); ok {
				log.Fatalln("KinesisLoop:", awsErr.Code(), awsErr.Message())
			} else {
				log.Fatalln("KinesisLoop:", err)
			}
		} else if *resp.StreamDescription.StreamStatus == "ACTIVE" {
			break
		} else {
			log.Println("KinesisLoop:", KinesisStreamName, "not active")
			time.Sleep(1 * time.Second)
		}
	}

	for one_event := range event {
		msg := callback(one_event)
		if data, err := json.Marshal(msg); err != nil {
			log.Fatalln(err)
		} else {
			log.Println("KinesisLoop:", string(data))
			// TODO: send Kinesis event
		}
	}
}
