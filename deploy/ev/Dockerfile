FROM alpine:latest

WORKDIR /home/service/ev

COPY bin/ev .

COPY config/config.yml .

COPY deploy/ev/run.sh .

ENV GOPROXY=https://goproxy.cn,direct

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories && \
    apk add -U tzdata && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime  && \
    mkdir ./file

EXPOSE 8090

RUN chmod +x run.sh
ENTRYPOINT  ./run.sh
