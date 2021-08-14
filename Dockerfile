FROM golang:1.15-alpine AS build
ARG GO_OS="linux"
ARG GO_ARCH="amd64"
WORKDIR /build/
COPY . .

# Build binary output
RUN GOPROXY=https://goproxy.cn,direct GOOS=${GO_OS} GOARCH=${GO_ARCH} go get -u github.com/swaggo/swag/cmd/swag@v1.6.7
RUN GOPROXY=https://goproxy.cn,direct GOOS=${GO_OS} GOARCH=${GO_ARCH} swag init
RUN GOPROXY=https://goproxy.cn,direct GOOS=${GO_OS} GOARCH=${GO_ARCH} go build -o sensiblequery -ldflags '-s -w' main.go

FROM alpine:latest
RUN apk add tzdata && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && echo "Asia/Shanghai" > /etc/timezone

RUN adduser -u 1000 -D sato -h /data
USER sato
WORKDIR /data/

COPY --chown=sato --from=build /build/sensiblequery /data/sensiblequery
COPY --chown=sato --from=build /build/docs /data/docs

ENV LISTEN 0.0.0.0:8000
EXPOSE 8000
CMD ["./sensiblequery"]
