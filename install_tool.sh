#!/usr/bin/env bash

set -e
sudo apt update
# Assuming that GCC has been preinstalled on most Linux distribution, I didn't include it here
sudo apt install golang-go
go get -u github.com/davidbyttow/govips/v2/vips
