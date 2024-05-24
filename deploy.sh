#!/bin/sh

bazel run //k8s/rest-api:chart.upgrade

bazel run //src/payload-processor/cmd:push -- --tag "payload-processor-$(git rev-parse --short HEAD)"
helm upgrade payload-processor --set image.tag="payload-processor-$(git rev-parse --short HEAD)" --set db.password=$DB_PASSWORD k8s/payload-processor
