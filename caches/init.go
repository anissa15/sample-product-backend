package caches

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

type Redis struct {
	rc *redis.Client
}

func New(redisOptionsStr string) *Redis {
	var addr, password string
	var db, protocol int
	for _, v := range strings.Split(redisOptionsStr, " ") {
		s := strings.Split(v, "=")
		if len(s) < 2 {
			continue
		}
		switch s[0] {
		case "addr":
			addr = s[1]
		case "password":
			password = s[1]
		case "db":
			db, _ = strconv.Atoi(s[1])
		case "protocol":
			protocol, _ = strconv.Atoi(s[1])
		}
	}
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
		Protocol: protocol,
	})
	// test connection by set and get sample data
	ctx := context.Background()
	if err := client.Set(ctx, "ping", "pong", 1*time.Minute).Err(); err != nil {
		panic(err)
	}
	if _, err := client.Get(ctx, "ping").Result(); err != nil {
		panic(err)
	}
	return &Redis{rc: client}
}
