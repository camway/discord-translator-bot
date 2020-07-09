FROM golang:latest as builder

RUN apt-get update && \
    apt-get install -y inotify-tools

WORKDIR /go/src/translatorbot

RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

COPY ./entry-point /go/entry-point
RUN chmod +x /go/entry-point

ENTRYPOINT ["/go/entry-point"]
