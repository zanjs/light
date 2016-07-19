#!/bin/bash

set -e

# Install gobatis
go-bindata template.txt
go install

# Use persist to generate code
go generate examples/mapper/model.go

# Test
go test -v examples/mapper/*.go
