package main

import (
	"net"
	"sync"

	log "github.com/Sirupsen/logrus"
	"zombiezen.com/go/capnproto2/rpc"

	"docs.greenitglobe.com/despiegk/gotests/capnp_benchmarks/go/store"
)

//Server listens on a connection for incoming connections
type Server struct {
	laddr          string
	maxConnections int

	lis net.Listener

	// use a mutex to make the connections safe for access by multiple goroutines
	// by conventions, mutexes go above the variables they guard
	clientconnectionmutex sync.Mutex // protects following
	connections           []*ClientConnection
}

//NewServer creates a tcp server for listening on the local network address laddr.
// During the Start() call, a listening socket is created ( https://golang.org/pkg/net/#Listen ) using "tcp" as network and laddr as specified.
func NewServer(laddr string) (server *Server) {
	server = &Server{laddr: laddr, maxConnections: 1000}
	return
}

// ClientConnection maintains a connection to a client and (de)serializes requests/reponses/notifications
type ClientConnection struct {
	server *Server

	socket net.Conn

	User string
}

//NewClientConnection creates a new ClientConnection given a socket
func (server *Server) NewClientConnection(socket net.Conn) (c *ClientConnection) {
	return &ClientConnection{socket: socket, server: server}
}

//Start creates  connections on the listener and serves requests for each incoming connection.
// Start blocks until the underlying tcp listener returns a non-nil error or Close is called on the server.
// Start returns an `err error`, this is a named return value. the err variable will
// be initialized for us (so we don't have to do it ourselved), and a `naked return`
// (just a return statement) will return err. (https://tour.golang.org/basics/7)
func (server *Server) Start() (err error) {
	// declare an anonymous function. since anonymous functions are declared inline
	// in other function bodys, they can access the surrounding variables.
	// also see https://gobyexample.com/closures
	// this function will set up our server
	func() {
		// access the members of the server passed to Start()
		// lock the mutex
		server.clientconnectionmutex.Lock()
		// defer the unlock of said mutex so we are sure it always gets called
		defer server.clientconnectionmutex.Unlock()
		server.lis, err = net.Listen("tcp", server.laddr)
		// make a new array to hold the *ClientConnections, length is 0 since we don't have any
		// and don't want a faulty initialization, capacity is 10 so we have an initial capacity
		// specifying capacity is optional
		server.connections = make([]*ClientConnection, 0, 10)
		log.Infoln("Listening for incoming connections on", server.laddr)
		// after declaring the function, immediatly invoke it
	}()
	// err should always be nil (the default of an error type) since it wasn't assigned yet
	if err != nil {
		return
	}
	// for {} is an infinite loop
	for {
		// error is the result of another anonymous function, also with a named return
		err = func() (err error) {
			// Accept blocks until a connection comes in we can actually try to accept
			conn, err := server.lis.Accept()
			// check if there was an error accepting the connection
			if err != nil {
				// naked return, equal to return err
				return
			}
			// the connection is accepted, lock the mutex
			server.clientconnectionmutex.Lock()
			// don't forget to unlock the mutex
			defer server.clientconnectionmutex.Unlock()
			// create a new ClientConnection to store some info about this connection
			c := server.NewClientConnection(conn)
			// if we already have the max amount of connections, close our new connection
			// and print an error log message
			if len(server.connections) >= server.maxConnections {
				log.Errorln("Maximum number of client connections reached (", server.maxConnections, "), dropping connection request")
				c.Close()
			}
			// store the new connection in the server.
			server.connections = append(server.connections, c)
			//TODO: release the closed clientconnections

			// log a message that we have a new tcp connection from the remote, but only
			// if we enabled debug logging
			log.Debugln("New tcp connection from ", conn.RemoteAddr())
			// go c.Listen() invokes the Listen() function on c in a new goroutine. the current code
			// will move forward to the next statement. see https://tour.golang.org/concurrency/1
			// note that we can't listen for any possible return values, should there be any
			go c.Listen()
			// return the named variable, again this is equivalent to return err
			return
		}()
		// check if there is an actual error and break out of the loop
		if err != nil {
			// return the named variable, return == return err
			return
		}
	}
}

//Close releases the underlying tcp listener
func (server *Server) Close() {
	if server.lis != nil {
		server.lis.Close()
	}
}

//Close releases the tcp connection
func (c *ClientConnection) Close() {
	if c.socket != nil {
		c.socket.Close()
		log.Debugln("Closed tcp connection from", c.socket.RemoteAddr())
	}
}

//Listen reads data from the open connection, deserializes it and handles the reponses and notifications
// This is a blocking function and will continue to listen until an error occurs (io or deserialization)
func (c *ClientConnection) Listen() {
	// Create a new locally implemented StoreFactory
	s := store.StoreFactory_ServerToClient(storeFactory{})
	// Listen for calls, using the StoreFactory as the bootstrap interface.
	conn := rpc.NewConn(rpc.StreamTransport(c.socket), rpc.MainInterface(s.Client))
	// Wait for connection to abort.
	// Wait blocks until the remote closes the connection.
	conn.Wait()
	// the remote has closed the connection, now we close the connection we created in
	// our ClientConnection
	c.Close()
}
