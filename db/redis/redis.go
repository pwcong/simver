package redis

import (
	"strconv"

	"github.com/go-redis/redis"
)

var Client *redis.Client

func Open(ip string, port int, password string, db int) error {

	Client = redis.NewClient(&redis.Options{
		Addr:     ip + ":" + strconv.Itoa(port),
		Password: password,
		DB:       db,
	})

	_, err := Client.Ping().Result()

	return err

}

func Close() error {
	return Client.Close()
}
