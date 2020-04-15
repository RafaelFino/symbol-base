package redis

import (
	"time"

	"github.com/go-redis/redis"
)

type Connection struct {
	conn *redis.Client
}

type Config struct {
	Address  string
	Password string
	DB       int
}

func New(cfg *Config) (*Connection, error) {
	conn := redis.NewClient(&redis.Options{
		Addr:     cfg.Address,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	_, err := conn.Ping().Result()

	return &Connection{conn: conn}, err

}

func (c *Connection) SendCommand(args ...interface{}) (interface{}, error) {
	return c.conn.Do(args...).Result()
}

func (c *Connection) Get(key string) (string, error) {
	return c.conn.Get(key).Result()
}

func (c *Connection) Set(key string, value interface{}, expiration time.Duration) error {
	return c.conn.Set(key, value, expiration).Err()
}
