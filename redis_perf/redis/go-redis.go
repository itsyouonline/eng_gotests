package redis

import (
	go_redis "github.com/go-redis/redis"
)

// GoRedisClient is a wrapper for a "go-redis" client
type GoRedisClient struct {
	client  *go_redis.Client
	network string
}

// newGoRedisClient creates a new go redis client
func newGoRedisClient(connectionAddr string, conType ConnectionType) RedisClient {
	network := connTypeToString(conType)

	client := go_redis.NewClient(&go_redis.Options{
		Network:  network,
		Addr:     connectionAddr,
		Password: "",
		DB:       0,
	})
	return GoRedisClient{
		client:  client,
		network: network,
	}
}

func (rc GoRedisClient) String() string {
	return "go-redis - connection: " + rc.network
}

func (rc GoRedisClient) Ping() error {
	return rc.client.Ping().Err()
}

func (rc GoRedisClient) StoreInHset(key, field string, value []byte) error {
	bcmd := rc.client.HSet(key, field, value)
	if bcmd.Err() != nil {
		return bcmd.Err()
	}
	return nil
}

func (rc GoRedisClient) GetFromHset(key, field string) ([]byte, error) {
	cmd := rc.client.HGet(key, field)
	if cmd.Err() != nil {
		return nil, cmd.Err()
	}
	return []byte(cmd.Val()), nil
}
