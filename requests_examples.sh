#!/bin/sh

curl -H "API_KEY: ${FS_TRACER_API_KEY}" -X POST -d '
{
	"timestamp": "2020-01-02T15:04:05Z",
	"absolute_path": "/home/user/file.txt",
	"contents": "Hello, World!"
}
' http://leunam.dev:9999/file/

curl -H "API_KEY: ${FS_TRACER_API_KEY}" -X GET http://leunam.dev:9999/file/?path=%2Fhome%2Fuser%2Ffile.txt

