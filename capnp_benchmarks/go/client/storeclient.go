package main

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
	socket net.Conn
	conn   *rpc.Conn
}

//Dial connects to a stratum+tcp at the specified network address.
// This function is not threadsafe
func (c *Client) Dial(host string) (err error) {
	c.socket, err = net.Dial("tcp", host)
	c.conn = rpc.NewConn(rpc.StreamTransport(c.socket))
	return
}

//Close releases the tcp connection
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
	ctx := context.Background()

	// Get the "bootstrap" interface.  This is the capability set with
	// rpc.MainInterface on the remote side.
	sf := store.StoreFactory{Client: c.conn.Bootstrap(ctx)}

	s := sf.CreateStore(ctx, func(p store.StoreFactory_createStore_Params) error {
		jwt, e := p.NewJwt()
		if e != nil {
			return e
		}
		e = jwt.SetPayload([]byte("ROB"))
		return e
	}).Store()

	for _, iterations := range []int{1000, 10000, 50000} {
		promises := make([]store.Store_get_Results_Promise, iterations, iterations)
		startTime := time.Now()
		for i := 0; i < iterations; i++ {
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
			promises[i] = s.Get(ctx, func(p store.Store_get_Params) error {
				return nil
			})

			if err != nil {
				return
			}
		}
		for _, promise := range promises {
			promise.Struct()
		}
		log.Infoln(iterations, "times get/set took ", time.Since(startTime).Seconds(), "seconds")

	}
	return nil
}
