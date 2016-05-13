#!/bin/bash

set -e

# install persist
go-bindata persist.txt
go install

# use persist to generate code
go generate ./...

# test
go test -v examples/persist/*.go
