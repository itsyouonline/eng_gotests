package main

import (
	"os"

	"docs.greenitglobe.com/despiegk/gotests/encryption_test/perf"

	log "github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
)

var version string

func main() {
	if version == "" {
		version = "Dev"
	}
	app := cli.NewApp()
	app.Name = "Encryption and compression test"
	app.Version = version

	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	log.SetOutput(os.Stdout)

	var debugLogging bool

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:        "debug, d",
			Usage:       "Enable debug logging",
			Destination: &debugLogging,
		},
	}

	app.Before = func(c *cli.Context) error {
		if debugLogging {
			log.SetLevel(log.DebugLevel)
			log.Debug("Debug logging enabled")
		}
		return nil
	}

	app.Action = func(c *cli.Context) error {
		log.Infoln(app.Name, "version", app.Version)
		err := perf.TestIndividual(200)
		if err != nil {
			return err
		}
		return perf.TestList(200)
	}

	app.Run(os.Args)
}
