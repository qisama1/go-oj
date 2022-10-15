package utils

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

var ctx = context.Background()
var rdb = redis.NewClient(&redis.Options{
	Addr:     "127.0.0.1:6379",
	Password: "123456", // no password set
	DB:       0,        // use default DB
})

func Set(k, v string) error {
	err := rdb.Set(ctx, k, v, time.Second*60)
	return err.Err()
}

func Get(k string) (string, error) {
	v, err := rdb.Get(ctx, k).Result()
	return v, err
}
