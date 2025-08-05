#!/usr/bin/env bash

set -e
cat samples.txt | sort | uniq -c | dlv --headless debug main.go -- "$@"
