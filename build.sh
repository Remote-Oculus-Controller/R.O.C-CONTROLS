#!/usr/bin/env bash

go build -ldflags "-X /roll/core.Build=`git rev-parse HEAD`" -o roll main.go
