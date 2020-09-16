# syntax = docker/dockerfile:experimental
FROM golang:1.14.0
RUN apt-get update && apt-get install unzip
COPY . /go/src/github.com/grafeas/grafeas/
WORKDIR /go/src/github.com/grafeas/grafeas
RUN make build
WORKDIR /go/src/github.com/grafeas/grafeas/go/v1beta1/main
RUN --mount=type=cache,target=/root/.cache/go-build GO111MODULE=on CGO_ENABLED=0 go build -o grafeas-server .

FROM alpine:latest
WORKDIR /
COPY --from=0 /go/src/github.com/grafeas/grafeas/go/v1beta1/main/grafeas-server /grafeas-server
EXPOSE 8080
ENTRYPOINT ["/grafeas-server"]
