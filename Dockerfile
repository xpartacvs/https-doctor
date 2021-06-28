FROM golang:1.16-alpine AS builder
WORKDIR /builder/src
COPY . .
RUN mkdir -p /builder/bin
RUN go build -ldflags="-s -w" -o /builder/bin/https-doctor main.go

FROM alpine:latest
LABEL maintainer="xpartacvs@gmail.com"
WORKDIR /usr/local/bin
RUN apk update
RUN apk add --no-cache tzdata
COPY --from=builder /builder/bin/https-doctor .
CMD ["https-doctor"]
