package main

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	disc "github.com/jeffjen/agent/discovery"
	d2k "github.com/jeffjen/go-dockerevents"
	"path"
	"time"
)

const (
	INVALIDATE_SERVICE_TTL = 1 * time.Second
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
				switch ev.Status {
				case "start":
					log.WithFields(log.Fields{"discovery": dflag, "addr": addr, "ttl": ttl}).Infof("Registering instance: %s", ev.ID)
					err = dkv.Set(path.Join(INSTANCE_ROOT, ev.ID), string(payload), ttl)
					break
				case "stop":
					log.WithFields(log.Fields{"discovery": dflag, "addr": addr, "ttl": INVALIDATE_SERVICE_TTL}).Infof("Invalidate instance: %s", ev.ID)
					err = dkv.Set(path.Join(INSTANCE_ROOT, ev.ID), "", INVALIDATE_SERVICE_TTL)
					break
				}
				if err != nil {
					log.Warning(err)
					time.Sleep(1 * time.Second)
					continue
				}
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
		case "stop":
			log.Infof("Spot invalid instance: %s", event.ID)
			etcdpub <- event
			break
		}
	}
}
