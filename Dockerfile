FROM golang:1.18-alpine AS builder
WORKDIR /build
RUN apk update && apk upgrade && apk add git
COPY . /build/
RUN go build
    
FROM alpine:latest
COPY --from=builder /build/multilog_exporter /multilog_exporter
CMD ["sh", "-c", "/multilog_exporter --metrics.listen-addr ${MLEX_LISTEN:-0.0.0.0:9144} --config.file /mlex.yaml"]
