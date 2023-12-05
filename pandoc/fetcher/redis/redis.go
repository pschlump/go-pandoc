package redis

import (
	"fmt"
	"os"

	"github.com/pschlump/dbgo"
	"github.com/pschlump/go-pandoc/config"
	"github.com/pschlump/go-pandoc/pandoc/fetcher"
)

type RedisFetcher struct {
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
}

func NewRedisFetcher(conf config.Configuration) (dataFetcher fetcher.Fetcher, err error) {
	dataFetcher = &RedisFetcher{}

	// get config stuff
	redisConnectionString := conf.GetString("connect", "bad-missing-connection-string")
	redisAuth := conf.GetString("auth", "")
	dbgo.Printf("%(LF)%(green) Redis Connection String ->%s<-\n", redisConnectionString)
	dbgo.Printf("%(LF)%(green) Redis Auth ->%s<-\n", redisAuth)

	// xyzzy - connect to redis

	return
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

	data = params.RedisKey

	return
}
