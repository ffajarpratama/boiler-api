package redis

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ffajarpratama/boiler-api/config"
	"github.com/redis/go-redis/v9"
)

type IFaceRedis interface {
	Set(key string, data interface{}, duration time.Duration) error
	Get(key string) (string, error)
	Del(key string) error
}

type Redis struct {
	client *redis.Client
}

func NewRedisClient(cnf *config.Config) (IFaceRedis, error) {
	conn := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cnf.Redis.Host, cnf.Redis.Port),
		Password: cnf.Redis.Password,
		DB:       cnf.Redis.Database,
	})

	_, err := conn.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal("[redis-connection-error] \n", err.Error())
		return nil, err
	}

	return &Redis{client: conn}, nil
}

// Set implements IFaceRedis.
func (r *Redis) Set(key string, data interface{}, duration time.Duration) error {
	return r.client.Set(context.Background(), key, data, duration).Err()
}

// Get implements IFaceRedis.
func (r *Redis) Get(key string) (string, error) {
	return r.client.Get(context.Background(), key).Result()
}

// Del implements IFaceRedis.
func (r *Redis) Del(key string) error {
	_, err := r.Get(key)
	if err != nil {
		return err
	}

	return r.client.Del(context.TODO(), key).Err()
}
