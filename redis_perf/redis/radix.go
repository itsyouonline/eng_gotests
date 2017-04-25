package redis

import (
	log "github.com/Sirupsen/logrus"

	"github.com/mediocregopher/radix.v2/redis"
)

// RedigoClient is a wrapper for a "redigo" client. It is not thread safe
type RadixClient struct {
	client  *redis.Client
	network string
}

type RadixPipe struct {
	client *redis.Client
}

// newRadixClient creates a new radix client
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

func (rp RadixPipe) StoreInHset(key, field string, value []byte) error {
	rp.client.PipeAppend("HSET", key, field, value)
	return nil
}

func (rc RadixClient) GetFromHset(key, field string) ([]byte, error) {
	resp, err := rc.client.Cmd("HGET", key, field).Bytes()
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (rp RadixPipe) GetFromHset(key, field string) ([]byte, error) {
	rp.client.PipeAppend("HGET", key, field)
	return nil, nil
}

func (rp RadixPipe) Execute() ([]byte, error) {
	// Execute the pipe, check if there are errors
	var err error
	for err == nil {
		err = rp.client.PipeResp().Err
	}
	// ignore pipelineempty errors
	if err == redis.ErrPipelineEmpty {
		err = nil
	}
	return nil, err
}

func (rc RadixClient) StartPipe() RedisPipe {
	return RadixPipe{
		client: rc.client,
	}
}
