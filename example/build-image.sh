#!/bin/zsh

env GOOS=linux GOARCH=amd64 go build -o example-linux-amd64 example.go

docker build -t gregorpirolt/ui:latest .

docker run -it gregorpirolt/ui:latest
