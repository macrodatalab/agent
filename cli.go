package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"os"
	"path"
)

func main() {
	app := cli.NewApp()
	app.Name = path.Base(os.Args[0])
	app.Usage = "A monitor tool for bigobject instances"

	app.Author = "YI-HUNG JEN"
	app.Email = "yihungjen@macrodatalab.com"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "log-level, l",
			Value: "info",
			Usage: fmt.Sprintf("Log level (options: debug, info, warn, error, fatal, panic)"),
		},
	}

	// logs
	app.Before = func(c *cli.Context) error {
		log.SetOutput(os.Stderr)
		level, err := log.ParseLevel(c.String("log-level"))
		if err != nil {
			log.Fatalf(err.Error())
		}
		log.SetLevel(level)

		return nil
	}

	app.Commands = commands

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
