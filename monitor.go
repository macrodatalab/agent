package main

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	disc "github.com/yihungjen/agent/discovery"
	d2k "github.com/yihungjen/go-dockerevents"
	"path"
	"time"
)

func publish(dflag string, addr string, ttl time.Duration) (event chan *d2k.Event) {
	event = make(chan *d2k.Event, 100)
	go func() {
		// init docker client for container inspection
		client, err := d2k.NewClient()
		if err != nil {
			log.Fatal("unable to establish docker client: %v", err)
		}
		for ev := range event {
			container, err := client.InspectContainer(ev.ID)
			if err != nil {
				log.Error(err)
				continue
			}
			portmap := container.NetworkSettings.PortMappingAPI()
			for idx, _ := range portmap {
				if portmap[idx].PrivatePort != 0 {
					portmap[idx].IP = addr
				}
			}
			payload, _ := json.Marshal(portmap)
			for {
				dkv, err := disc.New(dflag, "")
				if err != nil {
					log.Fatal("unable to establish connection to discovery -- %v", err)
				}
				err = dkv.Set(path.Join(INSTANCE_ROOT, ev.ID), string(payload), ttl)
				if err != nil {
					log.Warning(err)
					time.Sleep(1 * time.Second)
					continue
				}
				log.WithFields(log.Fields{"discovery": dflag, "addr": addr, "ttl": ttl}).Infof("Registering instance: %s", ev.ID)
				break
			}
		}
	}()
	return
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

	filter := d2k.ParseEventFilter(c.String("filter"))

	// init docker event loop
	eventSink := d2k.EventLoop(filter, 100)

	// init docker event publisher
	etcdpub := publish(dflag, addr, ttl)

	log.Println("begin monitor process...")
	for event := range eventSink {
		switch event.Status {
		case "start":
			log.Infof("Obtained instance: %s", event.ID)
			etcdpub <- event
			break
		}
	}
}
