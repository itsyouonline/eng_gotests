package redis

type RedisPipe interface {
	// StoreInHset stores value in an hset with the designated key in the designated field
	StoreInHset(key, field string, value []byte) error
	// GetFromHset gets a value fron an hset with the designated key from the designated field
	GetFromHset(key, field string) ([]byte, error)
	// Executes the pipe and gets the results.
	Execute() ([]byte, error)
}
