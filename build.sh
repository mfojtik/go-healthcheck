#!/bin/sh

echo "Building ./healthchk"
go build -o ./healthchk cmd/healthchk.go
