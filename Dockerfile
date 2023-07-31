# 使用官方 Golang 镜像作为基础镜像
FROM golang:1.20-alpine as builder

# 设置工作目录
WORKDIR /go/src/github.com/sxz799/checkIPRecord



# 将应用的代码复制到容器中
COPY ./ .


# 编译应用程序
RUN go env -w GO111MODULE=on \
    && go env -w GOPROXY=https://goproxy.cn,direct \
    && go env \
    && go mod tidy \
    && go build -ldflags="-s -w"  -o app .


FROM alpine:latest

WORKDIR /home

RUN apk --no-cache add bash

RUN apk update && apk add tzdata
RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
RUN echo "Asia/Shanghai" > /etc/timezone


COPY --from=0 /go/src/github.com/sxz799/checkIPRecord/app ./


# 运行应用程序
CMD ["./app"]