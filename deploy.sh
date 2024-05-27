#!/bin/sh

bazel run //k8s/rest-api:chart.upgrade
bazel run //k8s/payload-processor:chart.upgrade
