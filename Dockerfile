# Build the Go App
FROM golang:1.11-alpine AS go-builder

RUN apk add --no-cache git && \
  go get -u github.com/golang/dep/cmd/dep

ARG VERSION

COPY . /go/src/github.com/grahamar/casher
RUN cd /go/src/github.com/grahamar/casher && \
  dep ensure -vendor-only && \
  go build -o /go/bin/casher -v github.com/grahamar/casher

FROM golang:1.11-alpine

WORKDIR /usr/local/bin

COPY --from=go-builder /go/bin/casher /usr/local/bin/casher
