// package declaration, this file is in the `main` package
package main

// import other packages we want to use
import (
	// import the `os` package from the standard runtime
	"os"

	// import the `constants` package. The compiler looks for this package in:
	// $GOPATH/src/docs.greenitglobe.com/despiegk/gotests/capnp_benchmarks/go/constants
	"docs.greenitglobe.com/despiegk/gotests/capnp_benchmarks/go/constants"

	// import the `logrus` package. The compiler looks for this package in:
	// $GOPATH/src/github.com/Sirupsen/logrus
	// tell the compiler we want to refer to this package as `log`, rather than `logrus`
	log "github.com/Sirupsen/logrus"
	// import the `cli` package. The compiler looks for this package in:
	// $GOPATH/src/github.com/codegangsta/cli
	// we use this package to do the argument parsing for us
	"github.com/codegangsta/cli"
)

// main is the program entry point
func main() {
	// create a new app in the cli library
	app := cli.NewApp()
	// set the app name
	app.Name = "Capnp test client"
	// load the app version from the constants package
	app.Version = constants.Version
	if app.Version == "" {
		app.Version = "Dev"
	}

	// tell the loggin package to print a timestamp in front of our log messages
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	// also tell it to print to stdout
	log.SetOutput(os.Stdout)

	// declare some variables to store possible command line flags
	var debugLogging bool
	var serverAddress string

	// declare the possbile flags
	app.Flags = []cli.Flag{
		// declare a bool flag.
		cli.BoolFlag{
			// the flag is --debug or -d, both are valid
			Name: "debug, d",
			// tell the user what this flag does
			Usage: "Enable debug logging",
			// store the value in the debugLogging variable. Since this is a boolflag,
			// debugLogging will be true if the flag is set, false otherwise
			Destination: &debugLogging,
		},
		// also declare a string flag
		cli.StringFlag{
			// the flag is --server or -s, both are valid
			Name: "server, s",
			// tell the user what this flag does
			Usage: "Capnp rpc server address",
			// set a default value in case the flag does not get set
			Value: "localhost:" + constants.RPCPort,
			// store the flags value (or the default if the flag is not set) in the
			// serverAddress variable
			Destination: &serverAddress,
		},
	}

	// app.Before is a variable in the cli package we use. It takes any function with
	// signature `func(*cli.Context) error`. This allows us to do some setup logic
	// before we jump into our program.
	app.Before = func(c *cli.Context) error {
		// first, tell our logger to log the app name and version
		log.Infoln(app.Name, "-", app.Version)
		// now, if the user has set the debug flag, tell the log package to set the
		// log level to debug. This will make sure logs like log.debug() are also printed.
		// the default log level is info. For more info, see github.com/Sirupsen/logrus
		if debugLogging {
			// set the log level to debug
			log.SetLevel(log.DebugLevel)
			// and print a debug messge to inform the user debug logging has been enabled
			log.Debugln("Debug logging enabled")
		}
		// the function signature expects us to return something of type error. since we
		// didn't encounter an error, just return nil
		return nil
	}

	// app.Action is like app.Before, except that it doesn't return anything. Here we should
	// launch our application
	app.Action = func(c *cli.Context) {
		// declare a new storeClient value, and set it to a `Client` struct
		storeClient := &Client{}
		// call the Dial function of the Client and pass the serverAddress variable.
		// Store the returned error in a newly declared `err` variable
		err := storeClient.Dial(serverAddress)
		// check if the Dial function actually returned an error
		if err != nil {
			// if the err variable is not nil, there actually was an error in the Dial function,
			// so lets log the error. Error is a build in type defined by the following interface:
			// type error interface {
			//	Error() string
			// }
			// this means that any type which defines a function Error() string on that type
			// can be used as an error. the Error() function should return a string representation
			// of what went wrong. the runtime functions used to print to the terminal are able
			// to detect error types, and will call the Error() method themselves when we
			// ask them to print the error, so we don't have to do log.Errorln(err.Error())
			// since our loging package uses these functions under the hood, we can just pass
			// our error and it will take care of converting it to a string for us
			log.Errorln(err)
			// if there was an error, we also don't want to continue. since we didn't specify
			// any return values in the function signature, this `empty` return just causes
			// our function to stop executing
			return
		}
		// if we get here, Dial did not return an error and as such we have an open connection
		// to the server. we want to make sure this connection is always closed when we return
		// from this function, so we use defer. a defered statement is executed just before
		// the return statement. if we have multiple code paths wich return, defer always gets executed
		// multiple defers are possible, in this case they will be executed in reverse order of declaration:
		// e.g.
		// defer stmt1
		// defer stmt2
		// defer stmt3
		// return
		// will execute stmt3, then stmt2, then stmt1 and finally return from the function.
		// it is also possible to alter the return values in a defer, or defer (anonymous)
		// functions
		defer storeClient.Close()

		// Now that everything is set up, execute the benchmark. this can also give an error,
		// so we reuse our empty `err` variable. Since err is already declared, we can just
		// use `=` to assign the returned error to the variable, rather than using `:=` to both
		// declare and assign the variable
		err = storeClient.ExecuteBenchmark()
		// now check if an error occurred again
		if err != nil {
			// if there was an error, print it. there is no real need to return here
			// as we've reached the end of our function and will return anyway
			log.Errorln(err)
		}
	}

	// now that we told the app what to do, start it. pass os.Args (the command line params).
	// the cli package will parse the flags, invoke the app.Before function and then the
	// app.Action function
	app.Run(os.Args)
}
