#!/bin/sh

if [ -z "${FS_TRACER_API_KEY}" ]; then
	echo "FS_TRACER_API_KEY is not set"
	exit 1
fi

curl -H "API_KEY: ${FS_TRACER_API_KEY}" -X POST -d '
[{
	"timestamp": "2020-01-02T15:04:05Z",
	"absolute_path": "/home/user/file.txt",
	"contents": "Hello, World!"
}]
' http://leunam.dev:9999/api/v1/file/

curl -H "API_KEY: ${FS_TRACER_API_KEY}" -X GET http://leunam.dev:9999/api/v1/file/?path=%2Fhome%2Fuser%2Ffile.txt

