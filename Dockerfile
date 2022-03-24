FROM golang:1.17-alpine as builder
WORKDIR /root/go/src/github.com/einsier/ustc_melee_game
COPY . /root/go/src/github.com/einsier/ustc_melee_game
RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh
RUN go build -o hall-server run.go

FROM alpine:latest
# environment variable
ARG GSRPC_ADDR
ARG DBPROXY_ADDR
ARG ETCD_ADDR
ENV ENV_GSRPC_ADDR=$GSRPC_ADDR \
    ENV_DBPROXY_ADDR=$DBPROXY_ADDR \
    ENV_ETCD_ADDR=$ETCD_ADDR
WORKDIR  /root/go/src/github.com/einsier/ustc_melee_game
COPY --from=builder  /root/go/src/github.com/einsier/ustc_melee_game/hall-server .
EXPOSE 8000/tcp
EXPOSE 9000/tcp
ENTRYPOINT ./hall-server -gsRpcAddr $ENV_GSRPC_ADDR -dbProxyAddr $ENV_DBPROXY_ADDR -etcdAddr $ENV_ETCD_ADDR