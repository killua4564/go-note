FROM	golang:1.15-alpine3.12 AS develop

ENV		GO111MODULE=on
WORKDIR /go/src/github.com/killua4564/go-note
COPY	. .

RUN		apk update && apk add --no-cache git
RUN		go mod download
RUN		go build .


FROM	alpine:3.12

COPY 	--from=develop /go/src/github.com/killua4564/go-note/go-note .
RUN		chmod +x go-note
CMD		["./go-note"]
