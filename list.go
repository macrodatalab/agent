package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	docker "github.com/fsouza/go-dockerclient"
	"github.com/jeffjen/go-dockerevents"
	"strings"
)

func ParseImagesToList(images string) (imgs map[string]bool) {
	if images == "" {
		imgs = nil
		return
	} else {
		imgs = make(map[string]bool)
		for _, oneimg := range strings.Split(images, ",") {
			imgs[oneimg] = true
		}
		return
	}
}

func list(c *cli.Context) {
	client, drr := dockerevents.NewClient()
	if drr != nil {
		log.Fatalf("Unable to reach docker host: %v", drr)
	}

	containers, err := client.ListContainers(docker.ListContainersOptions{})
	if err != nil {
		log.Fatal(err)
	}

	images := ParseImagesToList(c.String("images"))
	for _, oneContainer := range containers {
		if images == nil {
			fmt.Println(oneContainer)
		} else {
			if _, ok := images[oneContainer.Image]; ok {
				fmt.Println(oneContainer)
			}
		}
	}
}
