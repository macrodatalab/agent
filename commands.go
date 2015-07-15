package main

import (
	"github.com/codegangsta/cli"
	"os"
)

var (
	commands = []cli.Command{
		{
			Name:  "list",
			Usage: "list bigobject instances",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "images",
					Usage: "Limit what containers are reported by image",
				},
			},
			Action: list,
		},
		{
			Name:  "monitor",
			Usage: "monitor bigobject instance life cycle on this node",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "addr",
					Usage:  "Public Address of the bigobject instance",
					EnvVar: "BO_HOST_IP",
				},
				cli.StringFlag{
					Name:   "ttl",
					Usage:  "Time to live for bigobject instance",
					EnvVar: "BO_INST_TTL",
				},
				cli.StringFlag{
					Name:   "filter",
					Usage:  "Docker event filter in YAML",
					EnvVar: "DOCKER_EVENT_FILTER",
				},
			},
			Action: monitor,
		},
		{
			Name:   "watch",
			Usage:  "watch instances' expiration on this node",
			Action: watch,
		},
	}
)

func getDiscovery(c *cli.Context) string {
	if len(c.Args()) == 1 {
		return c.Args()[0]
	}
	return os.Getenv("SWARM_DISCOVERY")
}
