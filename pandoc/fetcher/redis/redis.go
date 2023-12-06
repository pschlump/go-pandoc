package redis

import (
	"context"
	"fmt"
	"os"

	"github.com/pschlump/dbgo"
	"github.com/pschlump/go-pandoc/config"
	"github.com/pschlump/go-pandoc/pandoc/fetcher"
	"github.com/redis/go-redis/v9"
)

var ctx context.Context

type RedisFetcher struct {
	rdb     *redis.Client
	AuthKey string
}

type Params struct {
	RedisKey string `json:"rediskey"`
	AuthKey  string `json:"authkey"`
}

func (param *Params) Validation(key string) (err error) {
	if param.RedisKey == "" {
		err = fmt.Errorf("[fetcher-redis]: params of reiskkey is empty")
		dbgo.Printf("%(LF) error %s\n", err)
		return
	}
	dbgo.Printf("%(yellow) param.AuthKey ->%s<- v.s. key ->%s<-\n", param.AuthKey, key)
	if param.AuthKey != key {
		err = fmt.Errorf("[fetcher-redis]: params of authkey did not match - not authorized")
		dbgo.Printf("%(LF) error %s\n", err)
		return
	}
	return
}

func init() {
	if err := fetcher.RegisterFetcher("redis", NewRedisFetcher); err != nil {
		panic(err)
	}
	ctx = context.Background()
}

func NewRedisFetcher(conf config.Configuration) (dataFetcher fetcher.Fetcher, err error) {
	var me RedisFetcher

	// get config stuff
	redisHostPort := conf.GetString("connect", "bad-missing-connection-string")
	redisAuth := conf.GetString("auth", "")
	dbgo.Printf("%(LF)%(green) Redis Connection String ->%s<-\n", redisHostPort)
	dbgo.Printf("%(LF)%(green) Redis Auth ->%s<-\n", redisAuth)

	// connect to redis
	if redisAuth != "" {
		me.rdb = redis.NewClient(&redis.Options{
			Addr:     redisHostPort,
			Password: redisAuth,
			DB:       0, // 0 is default DB
		})
	} else {
		me.rdb = redis.NewClient(&redis.Options{
			Addr: redisHostPort,
			DB:   0,
		})
	}
	me.AuthKey = os.Getenv(conf.GetString("auth-key"))
	dbgo.Printf("%(LF)%(green)AuthKey = ->%s<-\n", me.AuthKey)

	return &me, nil
}

func (fetch *RedisFetcher) Fetch(fetchParams fetcher.FetchParams) (data []byte, err error) {

	dbgo.Fprintf(os.Stderr, "%(red)%(LF) - fetching using 'redis' protocal\n")

	var params Params

	if err = fetchParams.Unmarshal(&params); err != nil {
		dbgo.Printf("%(LF) error %s\n", err)
		return
	}

	if err = params.Validation(fetch.AuthKey); err != nil {
		dbgo.Printf("%(LF) error %s\n", err)
		return
	}

	// need to do a "get" from redis.
	x, err := fetch.rdb.Get(ctx, string(params.RedisKey)).Result()
	dbgo.Printf("%(LF)%(cyan) key ->%s<- data = --->>>%s<<<---, err=%s\n", params.RedisKey, x, err)
	if err != nil {
		dbgo.Printf("%(LF)%(red) Error %s\n", err)
		return
	}

	// data = params.RedisKey
	data = []byte(x)

	return
}
