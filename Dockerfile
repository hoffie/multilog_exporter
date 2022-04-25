FROM golang:1.18-alpine AS builder
COPY . /build
WORKDIR /build
RUN apk update && \
    apk upgrade && \
    apk add --no-cache git

RUN git clone https://github.com/kir4h/multilog_exporter.git && \
    cd multilog_exporter && \
    go mod init multilog && \
    go get ./... && \
    go mod vendor && \
    go build && \
    mv multilog ../ && \
    cd .. && \
    rm -fR multilog_exporter
    
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /build/multilog ./
CMD ["./multilog"]
