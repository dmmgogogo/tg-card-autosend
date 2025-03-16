#!/bin/bash

GOOS=linux GOARCH=amd64 go build -o tg-auto-card-num main.go
zip tg-auto-card-num.zip tg-auto-card-num