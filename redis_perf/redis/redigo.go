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

type RedigoPipe struct {
	conn redis.Conn
}

// newRedigoClient creates a new redigo client
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
	return "redigo - connection: " + rc.network
}

func (rc RedigoClient) Ping() error {
	_, err := rc.conn.Do("PING")
	return err
}

func (rc RedigoClient) StoreInHset(key, field string, value []byte) error {
	_, err := rc.conn.Do("HSET", key, field, value)
	return err
}

func (rp RedigoPipe) StoreInHset(key, field string, value []byte) error {
	return rp.conn.Send("HSET", key, field, value)
}

func (rc RedigoClient) GetFromHset(key, field string) ([]byte, error) {
	return redis.Bytes(rc.conn.Do("HGET", key, field))
}

func (rp RedigoPipe) GetFromHset(key, field string) ([]byte, error) {
	return nil, rp.conn.Send("HGET", key, field)
}

func (rp RedigoPipe) Execute() ([]byte, error) {
	// Execute the pipe, don't check the returns
	_, err := rp.conn.Do("")
	return nil, err
}

func (rc RedigoClient) StartPipe() RedisPipe {
	return RedigoPipe{
		conn: rc.conn,
	}
}
