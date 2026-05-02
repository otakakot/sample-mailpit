#!/bin/bash
set -e

rm -f swagger.json

go get -tool github.com/go-swagger/go-swagger/cmd/swagger@latest

curl -sL https://raw.githubusercontent.com/axllent/mailpit/master/server/ui/api/v1/swagger.json > swagger.json

mkdir -p ./pkg/mailpit/

go tool github.com/go-swagger/go-swagger/cmd/swagger generate client --quiet --spec=swagger.json --target=./pkg/mailpit/

go mod tidy
