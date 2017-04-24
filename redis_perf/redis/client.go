package redis

import (
	"fmt"
	"log"
)

// RedisClient is an interface defining all methods our custom redis client wrappers must implement
type RedisClient interface {
	// stringer interface is used to format the client when we print
	fmt.Stringer
	// Ping sends a ping to the redis server, this can be used to check the connection
	Ping() error
	// StoreInHset stores value in an hset with the designated key in the designated field
	StoreInHset(key, field string, value []byte) error
	// GetFromHset
	GetFromHset(key, field string) ([]byte, error)
}

type ConnectionType int

const (
	Tcp ConnectionType = iota
	Unix
)

// NewRedisClient creates a new RedisClient with the specified clientType. An unknown
// clientType is a fatal error
func NewRedisClient(clientType, connectionAddr string, conType ConnectionType) RedisClient {
	var client RedisClient
	switch clientType {
	case "go-redis":
		client = newGoRedisClient(connectionAddr, conType)
		break
	case "redigo":
		client = newRedigoClient(connectionAddr, conType)
		break
	case "radix":
		client = newRadixClient(connectionAddr, conType)
		break
	default:
		log.Fatal(clientType, " is not a recognized client.")
	}
	return client
}

// connTypeToString returns a string representation of a ConnectionType
func connTypeToString(conType ConnectionType) string {
	var network string
	switch conType {
	case Tcp:
		network = "tcp"
		break
	case Unix:
		network = "unix"
		break
	default:
		network = "tcp"
	}
	return network
}
