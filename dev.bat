@echo off
set GOOS=linux
go build main.go
set GOOS=windows
docker build -t scratch-go .
docker container run scratch-go