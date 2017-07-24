package main

import (
	"os"

	"github.com/itsyouonline/eng_gotests/redis_crypto_perf/perf"
	"github.com/itsyouonline/eng_gotests/redis_perf/redis"

	log "github.com/Sirupsen/logrus"

	"github.com/urfave/cli"
)

const (
	defaultObjectAmount = 1000000
	defaultObjectSize   = 200
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
	app.Name = "Redis crypto tests"
	app.Version = version

	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	log.SetOutput(os.Stdout)

	var debugLogging, lz4 bool
	var dbConnectionString, network, clientType string
	var objectAmount int
	var conType redis.ConnectionType

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:        "debug, d",
			Usage:       "Enable debug logging",
			Destination: &debugLogging,
		},
		cli.BoolFlag{
			Name:        "lz4",
			Usage:       "Use lz4 compression instead of gzip",
			Destination: &lz4,
		},
		cli.StringFlag{
			Name:        "connectionstring, c",
			Usage:       "Redis connection string",
			Value:       "localhost:6379",
			Destination: &dbConnectionString,
		},
		cli.StringFlag{
			Name:        "connectiontype, t",
			Usage:       "The type of connection to redis, either tcp or unix",
			Value:       "tcp",
			Destination: &network,
		},
		cli.StringFlag{
			Name:        "client",
			Usage:       "The underlying client to use in tests to connect to redis. \"go-redis\", \"redigo\" and \"radix\" are allowed",
			Value:       "redigo",
			Destination: &clientType,
		},
		cli.IntFlag{
			Name:        "data-size, s",
			Usage:       "The size of the data per object to be stored",
			Value:       200,
			Destination: &dataSize,
		},
		cli.IntFlag{
			Name:        "object-amount, a",
			Usage:       "The ammount of objects to be stored in redis",
			Value:       defaultObjectAmount,
			Destination: &objectAmount,
		},
	}

	app.Before = func(c *cli.Context) error {
		if debugLogging {
			log.SetLevel(log.DebugLevel)
			log.Debug("Debug logging enabled")
		}
		if objectAmount <= 0 {
			log.Debugf("Invalid object amount (%v), setting to default of %v", objectAmount, defaultObjectAmount)
			objectAmount = defaultObjectAmount
		}
		if dataSize <= 0 {
			log.Debugf("Invalid data size (%v), setting to default of %v", dataSize, defaultObjectSize)
			dataSize = defaultObjectSize
		}
		switch network {
		case "tcp":
			log.Debug("Clients will try to connect to redis using tcp")
			conType = redis.Tcp
			break
		case "unix":
			log.Debug("Clients will try to connect to redis using a unix socket")
			conType = redis.Unix
			break
		default:
			log.Fatal("Unrecognized connection type, only  \"tcp\" and \"unix\" are allowed")
		}
		return nil
	}

	app.Action = func(c *cli.Context) error {
		log.Infoln(app.Name, "version", app.Version)

		log.Debug("Connect to redis server at address ", dbConnectionString)

		client := redis.NewRedisClient(clientType, dbConnectionString, conType)
		return perf.CryptoTest(objectAmount, dataSize, lz4, client)
	}

	app.Run(os.Args)
}
