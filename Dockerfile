FROM	golang:1.15-alpine3.12

ENV		GO111MODULE=on

RUN		apk update && apk add git
RUN		go mod download
