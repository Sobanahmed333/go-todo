#stage 1, builder
FROM golang:alpine as builder
COPY . /app/

WORKDIR /app

ENV GO111MODULE=on
RUN cd /app && go mod download && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o socialvalidbackend .

#stage 2, app
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/ .

ENV WAIT_VERSION 2.7.2
ADD https://github.com/ufoscout/docker-compose-wait/releases/download/$WAIT_VERSION/wait /wait
RUN chmod +x /wait