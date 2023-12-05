#!/bin/bash

# A simple "data" based test that returns Hello text.

IP=127.0.0.1

curl -X POST \
  http://${IP}:9092/api/v1/convert \
  -H 'accept-encoding: gzip' \
  -H 'cache-control: no-cache' \
  -H 'content-type: application/json' \
  -d '{
	"fetcher": {
		"name": "data",
		"params": {
			"data": "IyMjIEhlbGxvCgo+IEdvLVBhbmRvYw=="
		}
	},
	"converter":{
		"from": "markdown",
	    "to":   "pdf",
	    "standalone": true,
	    "variable":{
	    	"CJKmainfont": "Vera",
	    	"mainfont":    "Vera",
	    	"sansfont":    "Vera",
	    	"geometry:margin":"1cm",
	    	"subject":"gsjbxx"
	    },
	    "template": "/Users/philip/go/src/github.com/pschlump/go-pandoc/data/docs.template"
	},
	"template": "binary"
}' --compressed -o out/test.pdf
