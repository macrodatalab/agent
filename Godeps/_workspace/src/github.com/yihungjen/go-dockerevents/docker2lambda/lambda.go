package main

import (
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/lambda"
	d2k "github.com/yihungjen/go-dockerevents"
	"log"
	"os"
	"time"
)

var (
	// AWS Lambda endpoint
	LambdaEndpoint string
)

func init() {
	LambdaEndpoint = os.Getenv("AWS_LAMBDA_ENDPOINT")
}

func LambdaLoop(event <-chan *d2k.Event, callback func(*d2k.Event) interface{}) {
	// main driver for pushing Lambda events to trigger Lambda Function

	svc := lambda.New(&aws.Config{Region: "us-west-2"})

	// check if the target Lambda function exist
	params := new(lambda.ListFunctionsInput)
	if resp, err := svc.ListFunctions(params); err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			log.Fatalln("LambdaLoop:", awsErr.Code(), awsErr.Message())
		} else {
			log.Fatalln("LambdaLoop:", err)
		}
	} else if resp.Functions == nil {
		log.Fatalln("LambdaLoop: enpoint not found:", LambdaEndpoint)
	} else {
		found := false
		for _, desc := range resp.Functions {
			found = found || (*desc.FunctionName == LambdaEndpoint)
		}
		if !found {
			log.Fatalln("LambdaLoop: enpoint not found:", LambdaEndpoint)
		}
	}

	// timer for triggering a Lambda invoke
	var ticker = time.Tick(10 * time.Second)

	var backlog []interface{}

	// Enpoint exist, beging polling for events
	for {
		select {
		case one_event, ok := <-event:
			if !ok {
				log.Fatalln("LambdaLoop: docker event channel closed")
				return
			}
			backlog = append(backlog, callback(one_event))
		case <-ticker:
			if len(backlog) != 0 {
				if data, err := json.Marshal(backlog); err != nil {
					log.Fatalln(err)
				} else {
					params := &lambda.InvokeInput{
						FunctionName: &LambdaEndpoint,
						Payload:      data,
					}
					if _, err := svc.Invoke(params); err != nil {
						log.Fatalln(err)
					}
					log.Println("LambdaLoop: Invoke function:", LambdaEndpoint, "Payload:", string(data))
				}
				backlog = backlog[:0]
			}
		}
	}
}
