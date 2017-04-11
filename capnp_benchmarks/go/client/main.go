package main

import (
	"os"

	"docs.greenitglobe.com/despiegk/gotests/capnp_benchmarks/go/constants"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
)

func main() {

	app := cli.NewApp()
	app.Name = "Capnp test client"
	app.Version = constants.Version
	if app.Version == "" {
		app.Version = "Dev"
	}

	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	log.SetOutput(os.Stdout)

	var debugLogging bool
	var serverAddress string

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:        "debug, d",
			Usage:       "Enable debug logging",
			Destination: &debugLogging,
		},
		cli.StringFlag{
			Name:        "server, s",
			Usage:       "Capnp rpc server address",
			Value:       "localhost:" + constants.RPCPort,
			Destination: &serverAddress,
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
		storeClient := &Client{}
		err := storeClient.Dial(serverAddress)
		if err != nil {
			log.Errorln(err)
			return
		}
		defer storeClient.Close()

		err = storeClient.ExecuteBenchmark()
		if err != nil {
			log.Errorln(err)
		}
	}

	app.Run(os.Args)
}
