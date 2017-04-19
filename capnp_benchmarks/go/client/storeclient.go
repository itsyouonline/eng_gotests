package main

// again, import all the packages we need in this file
import (
	"context"
	"net"
	"time"

	"docs.greenitglobe.com/despiegk/gotests/capnp_benchmarks/go/store"

	log "github.com/Sirupsen/logrus"
	"zombiezen.com/go/capnproto2/rpc"
)

// Client maintains a connection to the server and (de)serializes requests/reponses/notifications
type Client struct {
	// Our Client has 2 variables, a socket of type net.Conn and a pointer to a conn
	// of type *rpc.Conn. Note that these variables are not exported as their names start
	// with a lowercase letter. Only functions declared in the same package (main in this case)
	// can directly access the socket and conn variables of our client. Note that Client itself
	// starts with an uppercase letter and is thus exported.
	socket net.Conn
	conn   *rpc.Conn
}

//Dial connects to a stratum+tcp at the specified network address.
// This function is not threadsafe
// Dial is a function declared on a pointer to our `Client` type. To call the function, we first
// need to have a *Client. The client Dial is called on is then passed into the function
// with the variable name `c`. Since the Dial function is declared in the same package as
// the Client type, we can access the unexported variables inside Client in the function
func (c *Client) Dial(host string) (err error) {
	c.socket, err = net.Dial("tcp", host)
	c.conn = rpc.NewConn(rpc.StreamTransport(c.socket))
	return
}

//Close releases the tcp connection
// Simply check if the socket and conn of the given Client are defined, and call
// close of them if applicable
func (c *Client) Close() {
	if c.conn != nil {
		c.conn.Close()
	}
	if c.socket != nil {
		c.socket.Close()
	}
}

//ExecuteBenchmark performs the client testing code
func (c *Client) ExecuteBenchmark() (err error) {
	// call the Background() function from the context package, it just returns an empty
	// Context struct
	ctx := context.Background()

	// Get the "bootstrap" interface.  This is the capability set with
	// rpc.MainInterface on the remote side.
	sf := store.StoreFactory{Client: c.conn.Bootstrap(ctx)}

	// Call CreateStore on our storefactory. this functions takes another function as argument.
	s := sf.CreateStore(ctx, func(p store.StoreFactory_createStore_Params) error {
		// create a new jwt and check for errors
		jwt, e := p.NewJwt()
		if e != nil {
			// return the error if there is one
			return e
		}
		// no error so the jwt we got is valid, now try to set the payload of the jwt
		// and store a possible error
		e = jwt.SetPayload([]byte("ROB"))
		// return response of SetPayload, error or nil
		return e
		// call the Store() function on the result of sf.CreateStore(). Chaining function
		// calls is possible if functions only return 1 thing
	}).Store()

	// for loops are usually written like:
	// for index, value := range (slice/array) {
	//   loop body
	// }
	// note that we use := since we declare and assign new variables here
	// since we don't care about the index, the first variable is a _. This tells
	// the compiler that we know there should be something here, but we don't care
	// about it. []int{1000, 10000, 50000} creates a new slice
	for _, iterations := range []int{1000, 10000, 50000} {
		// promises is an array of type store.Store_get_Results_Promise, with length `iterations`
		// and capacity `iterations`
		promises := make([]store.Store_get_Results_Promise, iterations, iterations)
		// record the current time
		startTime := time.Now()
		// another kind of for loop, with a counter
		// declare and assign the counter i, loop while i is smaller than iterations,
		// and increment i after each loop
		for i := 0; i < iterations; i++ {
			// The capnp rpc logic, comes from the go file generated from the capnp schema
			// again, one of the parameters is a function
			s.Set(ctx, func(p store.Store_set_Params) error {
				o, e := p.NewObject()
				if e != nil {
					return e
				}
				o.SetNumber(uint32(i))
				e = o.SetTitle("Hello world")
				if e != nil {
					return e
				}
				return nil
			})
			// Store the get calls
			promises[i] = s.Get(ctx, func(p store.Store_get_Params) error {
				return nil
			})

			if err != nil {
				return
			}
		}
		// now that our loop is done, loopt over the storeed get calls again so we can
		// get a response. If we don't do this, the server will get a broken pipe error for
		// our get calls when it tries to write the response
		for _, promise := range promises {
			promise.Struct()
		}
		// print the amount of get/set calls and the time it took to complete them
		log.Infoln(iterations, "times get/set took ", time.Since(startTime).Seconds(), "seconds")

	}
	return nil
}
