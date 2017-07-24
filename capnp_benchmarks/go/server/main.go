// the server also has a main package
package main

// and some imports
import (
	"os"

	"github.com/itsyouonline/eng_gotests/capnp_benchmarks/go/constants"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
)

// the main() method, entrypoint for our application
func main() {
	// like the client we also use the cli library
	// do some cli setup again
	app := cli.NewApp()
	app.Name = "Capnp test server"
	app.Version = constants.Version
	if app.Version == "" {
		app.Version = "Dev"
	}

	// set the logger options for our log package
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	log.SetOutput(os.Stdout)

	// declare the variables to hold our command line flags
	// they will have their defalt value (false and "")
	var debugLogging bool
	var bindAddress string

	// definde the command line flags
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			// debug logging
			Name:        "debug, d",
			Usage:       "Enable debug logging",
			Destination: &debugLogging,
		},
		cli.StringFlag{
			// and the port to listen on
			// also set the default value to the one in the constants file should this not be set
			// Note the ":", this must be prepended to the port number. if the user sets
			// only the port number and omits the ":", the server will log an error about a
			// missing port and exit
			Name:        "bind, b",
			Usage:       "Capnp rpc bind address",
			Value:       ":" + constants.RPCPort,
			Destination: &bindAddress,
		},
	}

	// Same setup as the client: before the porgram really starts, print the name
	// and version and possibly enable debug logging
	app.Before = func(c *cli.Context) error {
		log.Infoln(app.Name, "-", app.Version)
		if debugLogging {
			log.SetLevel(log.DebugLevel)
			log.Debugln("Debug logging enabled")
		}
		return nil
	}

	// Start our application
	app.Action = func(c *cli.Context) {
		// Create a new server on the specified port.
		s := NewServer(bindAddress)
		// log the returned value of Start(). since this function only returns in case
		// there is an error, this log line only gets executed if an error occurs.
		// Errorln() is a function specific to our log package, it adds some extra info
		// to our log line (a red ERROR at the start if the terminal supports colors,
		// or an additional field that sais this is an error otherwise)
		log.Errorln(s.Start())
	}

	// Now that the setup is done, run our application and pass the command line flags/args
	app.Run(os.Args)
}
