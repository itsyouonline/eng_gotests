package main

import (
	"os"

	"docs.greenitglobe.com/despiegk/gotests/capnp_benchmarks/go/constants"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
)

func main() {

	app := cli.NewApp()
	app.Name = "Capnp test server"
	app.Version = constants.Version
	if app.Version == "" {
		app.Version = "Dev"
	}

	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	log.SetOutput(os.Stdout)

	var debugLogging bool
	var bindAddress string

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:        "debug, d",
			Usage:       "Enable debug logging",
			Destination: &debugLogging,
		},
		cli.StringFlag{
			Name:        "bind, b",
			Usage:       "Capnp rpc bind address",
			Value:       ":" + constants.RPCPort,
			Destination: &bindAddress,
		},
	}

	app.Before = func(c *cli.Context) error {
		log.Infoln(app.Name, "-", app.Version)
		if debugLogging {
			log.SetLevel(log.DebugLevel)
			log.Debugln("Debug logging enabled")
		}
		return nil
	}

	app.Action = func(c *cli.Context) {
		s := NewServer(bindAddress)
		log.Errorln(s.Start())
	}

	app.Run(os.Args)
}
