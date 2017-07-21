#!/bin/bash
version=$1
if [[ $version == "" ]]; then
    version="1.0_default"
fi
GOOS=linux GOARCH=amd64 go build -ldflags "-X main.Version=$version -X main.Buildstamp=`date -u '+%Y-%m-%d_%I:%M:%S%p'` -X main.Githash=`git rev-parse HEAD`" -v -o glb-client