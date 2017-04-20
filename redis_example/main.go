package main

import (
	"os"

	redis "github.com/go-redis/redis"

	log "github.com/Sirupsen/logrus"

	"github.com/urfave/cli"
)

var (
	version  string
	dataSize int
)

func main() {
	if version == "" {
		version = "Dev"
	}
	app := cli.NewApp()
	app.Name = "Redis capnp example"
	app.Version = version

	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	log.SetOutput(os.Stdout)

	var debugLogging bool
	var dbConnectionString string
	var messageAmount int

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:        "debug, d",
			Usage:       "Enable debug logging",
			Destination: &debugLogging,
		},
		cli.StringFlag{
			Name:        "connectionstring, c",
			Usage:       "Redis connection string",
			Value:       "localhost:6379",
			Destination: &dbConnectionString,
		},
		cli.IntFlag{
			Name:        "data-length, l",
			Usage:       "Length of additional data to be stored in the Tlog blocks",
			Value:       0,
			Destination: &dataSize,
		},
		cli.IntFlag{
			Name:        "message-amount, a",
			Usage:       "The ammount of capnp messages to be stored in redis",
			Value:       1000,
			Destination: &messageAmount,
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

		log.Debug("Connect to redis server at address ", dbConnectionString)

		client := redis.NewClient(&redis.Options{
			Addr:     dbConnectionString,
			Password: "", // no password set
			DB:       0,  // use default DB
		})

		pong, err := client.Ping().Result()
		if err != nil {
			log.Error("Error while pinging redis server: ", err)
			return err
		}
		log.Info("Ping response from redis server: ", pong)
		storeAndReadCapnpInHset(client, messageAmount)
		return nil
	}

	app.Run(os.Args)
}
