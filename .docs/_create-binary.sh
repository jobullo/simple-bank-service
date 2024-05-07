#!/bin/bash
cd ..
go mod verify
go mod download
cd cmd/http
swag init -parseDependency
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo .
cd ../..
