FROM golang:1.12.1-alpine as builder
WORKDIR /data
RUN apk add --no-cache make git
COPY . .
RUN make

FROM alpine:latest
RUN apk add --no-cache ca-certificates
WORKDIR /data
COPY --from=builder /data/mtfosbot /usr/bin
COPY config.default.yml config.yml
EXPOSE 10230
CMD ["/usr/bin/mtfosbot", "-f", "/data/config.yml", "-dbtool"]
