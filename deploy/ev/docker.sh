#!/bin/bash

# 启动容器挂在到本地 8090 端口

export CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GOPROXY=https://goproxy.cn

rm -f bin/ev

go build  -o  bin/ev cmd/ev/main.go

docker stop ev

docker rm ev

docker rmi -f ev

docker build -f deploy/ev/Dockerfile -t ev .

docker run --name ev --ulimit core=0 --log-opt max-size=500m --log-opt max-file=6  -p 8090:8090 -v /var/log/ev_logs:/home/service/ev/logs -d ev

rm -f bin/ev
