FROM golang:1.8

MAINTAINER Randy Cahyana "ranch@synoday.com"

RUN go get -u github.com/golang/dep/...

ADD . /go/src/github.com/synoday/service/user

WORKDIR /go/src/github.com/synoday/service/user

RUN dep ensure

WORKDIR /go/src/github.com/synoday/service/user/cmd/server

RUN go build -i -o usersvc

ENTRYPOINT ./usersvc