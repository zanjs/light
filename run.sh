#!/bin/bash

set -e

# install gobatis
go-bindata template.txt
go install

# use persist to generate code
go generate ./...

# test
go test -v examples/persist/*.go
