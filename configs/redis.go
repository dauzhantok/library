package configs

import (
	"fmt"
	"github.com/go-redis/redis/v7"
)

// Redis redis object
type Redis struct {
	Host     string `toml:"HOST" env:"REDIS_HOST"`
	Port     int    `toml:"PORT" env:"REDIS_PORT"`
	Password string `toml:"PASSWORD" env:"REDIS_PASSWORD"`
	Db       int    `toml:"DB" env:"REDIS_DB"`
}

// ConfigureRedis configures redis
func (r *Redis) ConfigureRedis() (*redis.Client, error) {
	address := r.GetConnectionString()
	redisConfig := &redis.Options{
		Addr:     address,
		Password: r.Password,
		DB:       r.Db,
	}
	client := redis.NewClient(redisConfig)
	_, err := client.Ping().Result()
	if err != nil {
		return nil, err
	}
	return client, nil
}

// GetConnectionString returns redis connection string
func (r *Redis) GetConnectionString() string {
	return fmt.Sprintf("%s:%d", r.Host, r.Port)
}
