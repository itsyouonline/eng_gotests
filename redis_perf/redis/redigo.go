package redis

import (
	log "github.com/Sirupsen/logrus"

	"github.com/garyburd/redigo/redis"
)

// RedigoClient is a wrapper for a "redigo" client. It is not thread safe
type RedigoClient struct {
	conn    redis.Conn
	network string
}

// newRedigoClient creates a new go redis client
func newRedigoClient(connectionAddr string, conType ConnectionType) RedisClient {
	network := connTypeToString(conType)

	conn, err := redis.Dial(network, connectionAddr)
	if err != nil {
		log.Fatal("Failed to create redigo client: ", err)
	}
	return RedigoClient{
		conn:    conn,
		network: network,
	}
}

func (rc RedigoClient) String() string {
	return "go-redis - connection: " + rc.network
}

func (rc RedigoClient) Ping() error {
	_, err := rc.conn.Do("PING")
	return err
}

func (rc RedigoClient) StoreInHset(key, field string, value []byte) error {
	_, err := rc.conn.Do("HSET", key, field, value)
	if err != nil {
		return err
	}
	return nil
}

func (rc RedigoClient) GetFromHset(key, field string) ([]byte, error) {
	return redis.Bytes(rc.conn.Do("HGET", key, field))
}
