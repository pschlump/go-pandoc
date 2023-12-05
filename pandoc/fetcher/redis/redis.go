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
	rdb *redis.Client
}

type Params struct {
	RedisKey []byte `json:"rediskey"`
}

func (p *Params) Validation() (err error) {
	if len(p.RedisKey) == 0 {
		err = fmt.Errorf("[fetcher-data]: params of reiskkey is empty")
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

	return &me, nil
}

func (p *RedisFetcher) Fetch(fetchParams fetcher.FetchParams) (data []byte, err error) {

	dbgo.Fprintf(os.Stderr, "%(red)%(LF) - fetching using 'redis' protocal\n")

	var params Params

	if err = fetchParams.Unmarshal(&params); err != nil {
		return
	}

	if err = params.Validation(); err != nil {
		return
	}

	// xyzzy - need to do a "get" from redis.
	x, err := p.rdb.Get(ctx, "a").Result()
	dbgo.Printf("%(LF)%(red) x = %s, err=%s\n", x, err)

	data = params.RedisKey

	return
}
