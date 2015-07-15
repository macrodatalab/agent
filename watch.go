package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	docker "github.com/fsouza/go-dockerclient"
	disc "github.com/yihungjen/agent/discovery"
	d2k "github.com/yihungjen/go-dockerevents"
	"path"
	"time"
)

func follow(dflag string) (event chan *disc.Event) {
	event = make(chan *disc.Event)
	go func() {
		defer close(event)
		var waitindex uint64 = 0
		for {
			dkv, err := disc.New(dflag, INSTANCE_ROOT)
			if err != nil {
				log.Fatal("unable to establish connection to discovery -- %v", err)
			}
			active, _ := dkv.WatchTree("", waitindex)
			for active := range active {
				waitindex = active.Node.ModifiedIndex + 1
				event <- active
			}
			log.Warning("unexpected disconnection from discovery: unable to follow node")
			time.Sleep(1 * time.Second)
		}
	}()
	return
}

func watch(c *cli.Context) {
	dflag := getDiscovery(c)
	if dflag == "" {
		log.Fatalf("discovery required publish instance info. See '%s monitor --help'.", c.App.Name)
	}

	etcdSink := follow(dflag)

	log.Println("begin watch process...")
	for {
		select {
		case active, ok := <-etcdSink:
			if !ok {
				log.Fatalln("etcd sink disconnected")
			}
			switch active.Action {
			case "expire":
				key := path.Base(active.Node.Key)
				log.Infof("Observed instance expire: %s", key)
				client, _ := d2k.NewClient()
				client.RemoveContainer(docker.RemoveContainerOptions{
					ID:    key,
					Force: true,
				})
				break
			}
		}
	}
}
