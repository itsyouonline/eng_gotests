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
func (server *Server) Start() (err error) {
	func() {
		server.clientconnectionmutex.Lock()
		defer server.clientconnectionmutex.Unlock()
		server.lis, err = net.Listen("tcp", server.laddr)
		server.connections = make([]*ClientConnection, 0, 10)
		log.Infoln("Listening for incoming connections on", server.laddr)
	}()
	if err != nil {
		return
	}
	for {
		err = func() (err error) {
			conn, err := server.lis.Accept()
			if err != nil {
				return
			}
			server.clientconnectionmutex.Lock()
			defer server.clientconnectionmutex.Unlock()
			c := server.NewClientConnection(conn)
			if len(server.connections) >= server.maxConnections {
				log.Errorln("Maximum number of client connections reached (", server.maxConnections, "), dropping connection request")
				c.Close()
			}
			server.connections = append(server.connections, c)
			//TODO: release the closed clientconnections
			log.Debugln("New tcp connection from ", conn.RemoteAddr())
			go c.Listen()
			return
		}()
		if err != nil {
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
	conn.Wait()
	c.Close()
}
