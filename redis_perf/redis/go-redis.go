package redis

import (
	"strconv"
	"strings"

	go_redis "github.com/go-redis/redis"
)

// GoRedisClient is a wrapper for a "go-redis" client
type GoRedisClient struct {
	client  *go_redis.Client
	network string
}

type GoRedisPipe struct {
	pipe *go_redis.Pipeline
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
	return rc.client.HSet(key, field, value).Err()
}

func (rp GoRedisPipe) StoreInHset(key, field string, value []byte) error {
	return rp.pipe.HSet(key, field, value).Err()
}

func (rc GoRedisClient) GetFromHset(key, field string) ([]byte, error) {
	cmd := rc.client.HGet(key, field)
	if cmd.Err() != nil {
		return nil, cmd.Err()
	}
	return []byte(cmd.Val()), nil
}

func (rp GoRedisPipe) GetFromHset(key, field string) ([]byte, error) {
	cmd := rp.pipe.HGet(key, field)
	if cmd.Err() != nil {
		return nil, cmd.Err()
	}
	return []byte(cmd.Val()), nil
}

func (rc GoRedisClient) GetMemUsage() (int, error) {
	cmd := rc.client.Info("memory")
	if cmd.Err() != nil {
		return 0, cmd.Err()
	}
	fields := strings.Fields(cmd.String())
	for i := range fields {
		if strings.HasPrefix(fields[i], "used_memory:") {
			return strconv.Atoi(strings.TrimPrefix(fields[i], "used_memory:"))
		}
	}
	return 0, nil
}

func (rp GoRedisPipe) Execute() ([]byte, error) {
	// Execute the pipe, don't check the returns, only possible errors
	_, err := rp.pipe.Exec()
	return nil, err
}

func (rc GoRedisClient) StartPipe() RedisPipe {
	return GoRedisPipe{
		pipe: rc.client.Pipeline(),
	}
}
