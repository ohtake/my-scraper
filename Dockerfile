FROM golang:1.13 AS builder

WORKDIR /go/src/github.com/mono0x/my-scraper

ADD go.mod go.sum Makefile ./
RUN make download

ADD . ./
RUN make build-linux

FROM alpine:latest
RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /go/src/github.com/mono0x/my-scraper/my-scraper.linux /app
CMD ["/app"]
