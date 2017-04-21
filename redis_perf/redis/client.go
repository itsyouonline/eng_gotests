package redis

import "fmt"

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
