#!/usr/bin/env bash

# Build Server
cd server || exit
go build -v ./...

# Build client
cd ../client || exit
go build -v ./...

# Build library
cd ../libs/merkletree || exit
go build -v ./...
