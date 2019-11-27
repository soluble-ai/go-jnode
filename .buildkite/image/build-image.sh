#!/bin/bash
set -e
docker build -t gcr.io/soluble-ci/go-jnode-ci:latest .
docker push gcr.io/soluble-ci/go-jnode-ci:latest
