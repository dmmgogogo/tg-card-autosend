#!/bin/bash

GOOS=linux GOARCH=amd64 go build -o tg-card-autosed main.go
zip tg-card-autosed.zip tg-card-autosed