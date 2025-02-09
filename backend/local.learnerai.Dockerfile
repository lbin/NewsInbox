FROM golang:1.22.3 AS build

WORKDIR /go/src/learnerai
COPY . /go/src/learnerai

# RUN go build -v -mod=vendor -o /bin/learnerai
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=1 GOPROXY=https://goproxy.cn \
    go build -v \
    -o /bin/learnerai
FROM learnerai:latest
EXPOSE 8099
WORKDIR /output

RUN sed -i 's/archive.ubuntu.com/mirrors.aliyun.com/g' /etc/apt/sources.list
RUN apt-get update \
    && apt-get install -y --no-install-recommends ca-certificates curl
RUN apt-get install tzdata && \
    ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone

COPY --from=build /bin/learnerai /output/
COPY conf/api.cloud.yaml /output/conf/api.yaml
RUN mkdir -p runtime/logs
CMD ["./learnerai"]
