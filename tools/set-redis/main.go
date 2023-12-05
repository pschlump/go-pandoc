package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/pschlump/dbgo"
	"github.com/redis/go-redis/v9"
)

var ctx context.Context
var rdb *redis.Client

func main() {

	ctx = context.Background()

	key := os.Args[1]
	fn := os.Args[2]
	buf, err := ioutil.ReadFile(fn)
	if err != nil {
		fmt.Printf("File read error: %s for %s\n", err, fn)
		os.Exit(1)
	}

	redisHostPort := "127.0.0.1:6379"
	redisAuth := ""

	// connect to redis
	if redisAuth != "" {
		rdb = redis.NewClient(&redis.Options{
			Addr:     redisHostPort,
			Password: redisAuth,
			DB:       0, // 0 is default DB
		})
	} else {
		rdb = redis.NewClient(&redis.Options{
			Addr: redisHostPort,
			DB:   0,
		})
	}

	dbgo.Printf("%(green)%s set to %s\n", fn, buf)

	rdb.Set(ctx, key, string(buf), 10*60*time.Minute)

}
