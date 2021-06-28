#!/bin/sh

go get

echo $PATH

reflex -r '(\.go$|go\.mod)' -s go run cmd/guardmech/main.go