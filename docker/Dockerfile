FROM alpine
MAINTAINER "4xx.me"

VOLUME /tmp
ADD ElasticViewLinux/ /

ENV TZ Asia/Shanghai

RUN chmod +x /ElasticViewLinux

CMD ["/ElasticViewLinux","-configFileDir=config","-configFileName=config","-configFileExt=yml"]