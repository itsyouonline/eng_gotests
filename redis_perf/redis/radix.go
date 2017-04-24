package redis

import (
	"log"

	"github.com/mediocregopher/radix.v2/redis"
)

// RedigoClient is a wrapper for a "redigo" client. It is not thread safe
type RadixClient struct {
	client  *redis.Client
	network string
}

func newRadixClient(connectionAddr string, conType ConnectionType) RedisClient {
	network := connTypeToString(conType)

	client, err := redis.Dial(network, connectionAddr)
	if err != nil {
		log.Fatal("Failed to create radix client: ", err)
	}
	return RadixClient{
		client:  client,
		network: network,
	}
}

func (rc RadixClient) String() string {
	return "radix.v2 - connection: " + rc.network
}

func (rc RadixClient) Ping() error {
	return rc.client.Cmd("PING").Err
}

func (rc RadixClient) StoreInHset(key, field string, value []byte) error {
	return rc.client.Cmd("HSET", key, field, value).Err
}

func (rc RadixClient) GetFromHset(key, field string) ([]byte, error) {
	resp, err := rc.client.Cmd("HGET", key, field).Bytes()
	if err != nil {
		return nil, err
	}
	return resp, nil
}
