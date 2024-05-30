#!/bin/sh

if [ -z "$DB_PASSWORD" ]; then
  echo "DB_PASSWORD is not set"
  exit 1
fi

bazel run //k8s/rest-api:chart.upgrade
bazel run //k8s/payload-processor:chart.upgrade
