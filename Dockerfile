FROM golang:1.18-alpine AS builder
ADD . /multilog_exporter
WORKDIR /multilog_exporter

RUN apk update && \
    apk upgrade && \
    apk add --no-cache git && \
    go build

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /multilog_exporter/multilog_exporter ./
CMD ["./multilog_exporter"]
