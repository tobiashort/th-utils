#!/usr/bin/env bash

set -e
go build
cat samples.txt | sort | uniq -c | ./uniqplot "$@"
