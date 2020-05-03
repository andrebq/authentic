package tcache

import (
	"time"

	"github.com/go-redis/redis/v7"
)

type (
	// Redis implements a TokenSet using a Redis compatible storage
	Redis struct {
		conn *redis.Client
	}
)

// NewRedisFromClient returns a new TokenSet using the provided redis client
func NewRedisFromClient(cli *redis.Client) (*Redis, error) {
	if err := cli.Ping().Err(); err != nil {
		return nil, err
	}
	return &Redis{conn: cli}, nil
}

// NewRedis returns a Redis token set using the redis server
// at the specific address.
//
// If you need more control over the redis client configuration, use
// NewRedisFromClient.
//
// Timeouts for reads/write are quite aggressive being 1sec for reads
// and 3 seconds for writes.
//
// DialTimeout is more relaxed being 10sec by default
func NewRedis(addr string) (*Redis, error) {
	return NewRedisFromClient(redis.NewClient(&redis.Options{
		Addr:         addr,
		DialTimeout:  time.Second * 10,
		ReadTimeout:  time.Second * 1,
		WriteTimeout: time.Second * 3,
	}))
}

// Contains implements TokenSet Contains
func (r *Redis) Contains(k string) (bool, error) {
	n, err := r.conn.Get(k).Int()
	return n == 1, err
}

// Add includes the key in the set of valid sessions
func (r *Redis) Add(k string, ttl time.Time) error {
	return r.conn.Set(k, int(1), ttl.Sub(time.Now())).Err()
}
