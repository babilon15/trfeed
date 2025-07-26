#!/bin/bash

GOARCH=amd64 GOAMD64=v2 go build -trimpath -ldflags="-s -w"
