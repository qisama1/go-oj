package test

import (
	"context"
	"github.com/go-redis/redis/v8"
	"testing"
)

var ctx = context.Background()
var rdb = redis.NewClient(&redis.Options{
	Addr:     "127.0.0.1:6379",
	Password: "123456", // no password set
	DB:       0,        // use default DB
})

func TestRedis(t *testing.T) {
	err := rdb.Set(ctx, "go-redis", "value", 0).Err()
	if err != nil {
		panic(err)
	}
}

func TestGet(t *testing.T) {
	v, _ := rdb.Get(ctx, "go-redis").Result()
	println(v)
}
