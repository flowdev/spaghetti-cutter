#!/bin/sh

GOOS=linux go build -o spaghetti-cutter.linux
GOOS=darwin go build -o spaghetti-cutter.macos
GOOS=windows go build -o spaghetti-cutter.windows
