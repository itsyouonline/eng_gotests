package main

import (
	"os"

	"github.com/itsyouonline/eng_gotests/encryption_test/perf"

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
	var dataSize, dataSet int

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:        "debug, d",
			Usage:       "Enable debug logging",
			Destination: &debugLogging,
		},
		cli.IntFlag{
			Name:        "size, s",
			Usage:       "Set the amount of random data to store in a block",
			Value:       200,
			Destination: &dataSize,
		},
		cli.IntFlag{
			Name:        "amount, a",
			Usage:       "The amount of blocks to generate in tests",
			Value:       1000000,
			Destination: &dataSet,
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
		err := perf.TestIndividual(dataSet, dataSize)
		if err != nil {
			return err
		}
		err = perf.TestCombined(dataSet, dataSize)
		if err != nil {
			return err
		}
		return perf.TestList(dataSet, dataSize)
	}

	app.Run(os.Args)
}
