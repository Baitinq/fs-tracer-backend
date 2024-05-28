#!/bin/sh

curl -H "API_KEY: ${API_KEY}" -X POST -d '{"timestamp": "2017-01-02T15:04:05Z"}' http://leunam.dev:9999/file/

curl -H "API_KEY: ${API_KEY}" -X GET http://leunam.dev:9999/file/
