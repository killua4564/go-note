FROM	golang:1.15-alpine3.12

RUN		apk update && apk add git
RUN		go get -u github.com/caarlos0/env \
                  github.com/gin-gonic/gin \ 
                  github.com/go-sql-driver/mysql
