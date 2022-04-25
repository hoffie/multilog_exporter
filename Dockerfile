#FROM ubuntu:latest AS builder
FROM golang:1.18-alpine AS builder
RUN apk update && \
    apk upgrade && \
    apk add git

RUN mkdir /build
WORKDIR /build

RUN git clone https://github.com/kir4h/multilog_exporter.git . 
RUN go mod init multilog_exporter
RUN go get ./...
RUN go mod vendor
RUN go build
    
FROM alpine:latest
#FROM golang:1.18-alpine
#FROM ubuntu:latest
#RUN apk --no-cache add ca-certificates
#WORKDIR /root/
COPY --from=builder /build/multilog_exporter /root/multilog_exporter
CMD ["/root/multilog_exporter"]
