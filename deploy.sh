#!/bin/sh

bazel run //src/rest-api/cmd:push -- --tag "rest-api-$(git rev-parse --short HEAD)"
helm upgrade rest-api --set image.tag="rest-api-$(git rev-parse --short HEAD)" k8s/rest-api

bazel run //src/payload-processor/cmd:push -- --tag "payload-processor-$(git rev-parse --short HEAD)"

helm upgrade payload-processor --set image.tag="payload-processor-$(git rev-parse --short HEAD)" --set db.password=$DB_PASSWORD k8s/payload-processor
