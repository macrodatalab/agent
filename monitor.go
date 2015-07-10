package main

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/docker/libkv/store"
	disc "github.com/yihungjen/agent/discovery"
	d2k "github.com/yihungjen/go-dockerevents"
	"os"
	"time"
)

const (
	INSTANCE_ROOT = "instances"
)

func getDiscovery(c *cli.Context) string {
	if len(c.Args()) == 1 {
		return c.Args()[0]
	}
	return os.Getenv("SWARM_DISCOVERY")
}

func monitor(c *cli.Context) {
	dflag := getDiscovery(c)
	if dflag == "" {
		log.Fatalf("discovery required publish instance info. See '%s monitor --help'.", c.App.Name)
	}

	addr := c.String("addr")
	if addr == "" {
		log.Fatal("missing mandatory --addr flag")
	}

	ttl, err := time.ParseDuration(c.String("ttl"))
	if err != nil {
		log.Fatal("invlaid --ttl: %v", err)
	}

	dkv, drr := disc.New(dflag, &disc.Options{
		INSTANCE_ROOT,
		&store.Config{EphemeralTTL: ttl},
	})
	if drr != nil {
		log.Fatal(drr)
	}

	evFilter := d2k.ParseEventFilter(c.String("filter"))

	eventSink := make(chan *d2k.Event, 100)
	go d2k.EventLoop(eventSink, evFilter)

	log.Println("begin monitor process...")
	for event := range eventSink {
		switch event.Status {
		default:
			break
		case "start":
			// TODO: send monitored instances to discovery service
			payload, _ := json.Marshal(event)
			log.WithFields(log.Fields{"discovery": dflag, "addr": addr, "ttl": ttl}).Infof("Registering instance: %s", event.ID)
			if err := dkv.Put(event.ID, payload, &store.WriteOptions{Ephemeral: true}); err != nil {
				log.Error(err)
			}
			break
		}
	}
}
