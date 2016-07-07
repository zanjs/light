#!/bin/bash

set -e

# install gobatis
go-bindata template.txt
# go install

# use persist to generate code
go generate examples/persist/demo.go

# test
go test -v examples/persist/*.go
