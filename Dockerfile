FROM golang:latest

ADD . /go/src/parser

WORKDIR /go/src/parser
