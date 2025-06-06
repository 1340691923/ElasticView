# 使用基础镜像，比如 alpine 或者 debian，alpine 体积较小
FROM alpine:latest

VOLUME /app/config
VOLUME /app/ev_store

# 设置工作目录
WORKDIR /app

RUN mkdir -p /app/ev_store

# 将本地的二进制文件和配置文件复制到容器中
COPY ev_linux_amd64 /app/
COPY config /app/config/

# 确保二进制文件有可执行权限
RUN chmod +x /app/ev_linux_amd64

EXPOSE 8090

# 定义容器启动时执行的命令
CMD ["./ev_linux_amd64", "-configFile=config/config.yml"]
