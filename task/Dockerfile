FROM golang:1.8

MAINTAINER Randy Cahyana "ranch@synoday.com"

RUN go get -u github.com/golang/dep/...

ADD . /go/src/github.com/synoday/service/task

WORKDIR /go/src/github.com/synoday/service/task

RUN dep ensure

RUN go install github.com/synoday/service/task

ENTRYPOINT /go/bin/task