#!/bin/bash

# A simple "data" based test that returns Hello text.

# Generate a key for Redis
# redisKey="IyMjIEhlbGxvCgo+IEdvLVBhbmRvYw=="
# redisKey="abc"

( cd tools/gen-uuid ; go build )
( cd tools/set-redis ; go build )
redisKey="go-pandoc:$(tools/gen-uuid/gen-uuid)"

# Set the value in reids fo rthe markdown data.
# - ttl of 10 min on key
./tools/set-redis/set-redis "${redisKey}" test/redis1.md


IP=127.0.0.1

mkdir -p out

curl -X POST \
  http://${IP}:9092/api/v1/convert \
  -H 'accept-encoding: gzip' \
  -H 'cache-control: no-cache' \
  -H 'content-type: application/json' \
  -d "{
	\"fetcher\": {
		\"name\": \"redis\",
		\"params\": {
			\"rediskey\": \"${redisKey}\"
		}
	},
	\"converter\":{
		\"from\": \"markdown\",
	    \"to\":   \"rtf\",
	    \"variable\":{
	    	\"CJKmainfont\": \"Liberation Sans\",
	    	\"mainfont\":    \"Liberation Sans\",
	    	\"sansfont\":    \"Liberation Sans\",
	    	\"pagestyle\":   \"empty\",
	    	\"geometry:margin\":\"1mm\",
	    	\"subject\":\"gsjbxx\"
	    },
	    \"template\": \"/Users/philip/go/src/github.com/pschlump/go-pandoc/data/docs.template\"
	},
	\"template\": \"binary\"
}" --compressed -o out/test-redis.rtf
