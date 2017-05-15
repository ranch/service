FROM golang:1.8

MAINTAINER Randy Cahyana "ranch@synoday.com"

RUN go get -u github.com/golang/dep/cmd/dep

ADD . /go/src/github.com/synoday/service/task

WORKDIR /go/src/github.com/synoday/service/task

RUN dep ensure

WORKDIR /go/src/github.com/synoday/service/task/cmd/server

RUN go build -i -o tasksvc

ENTRYPOINT ./tasksvc