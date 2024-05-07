#!/bin/bash
cd ..
go mod vendor
cd cmd/http
swag init -parseDependency
go run .
cd ../..
