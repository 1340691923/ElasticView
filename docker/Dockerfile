FROM mysql:5.7

#设置免密登录
ENV MYSQL_ALLOW_EMPTY_PASSWORD yes

COPY setup.sh /mysql/setup.sh

COPY es_view.sql /mysql/es_view.sql

COPY privileges.sql /mysql/privileges.sql

COPY ElasticViewLinux /usr/local/ElasticViewLinux

COPY config /usr/local/config

COPY config/config.json /usr/local/config/config.json

#设置容器启动时执行的命令
CMD ["sh", "/mysql/setup.sh"]
