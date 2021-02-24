FROM golang:latest
MAINTAINER "antgan@163.com"
WORKDIR /data/go/dryan
ADD . /data/go/dryan
RUN go build .
EXPOSE 9999
ENTRYPOINT ["./dryan"]