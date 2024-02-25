#!/usr/bin/env bash

# Build Server
cd server || exit
go test -v ./...

# Build client
cd ../client || exit
go test -v ./...

# Build library
cd ../libs/merkletree || exit
go test -v ./...
