#!/bin/bash

env GOOS=darwin GOARCH=amd64 go build -o ./bin/task_darwin ./cmd/task/task.go
env GOOS=linux GOARCH=amd64 go build -o ./bin/task_linux ./cmd/task/task.go
env GOOS=windows GOARCH=amd64 go build -o ./bin/task.exe ./cmd/task/task.go