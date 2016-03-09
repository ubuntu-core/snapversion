FROM golang:alpine

ADD . /go/src/github.com/ubuntu-core/snapversion

RUN apk add --no-cache git && \
  go get github.com/zenazn/goji && \
  go install github.com/ubuntu-core/snapversion/cmd/snapversion && \
  apk del git && \
  go clean -i github.com/zenazn/goji

ENTRYPOINT /go/bin/snapversion

EXPOSE 8000
