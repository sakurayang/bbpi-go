#!/bin/bash
go build -o dist/main main.go
docker build -f Dockerfile -t gerardyang/pi .
