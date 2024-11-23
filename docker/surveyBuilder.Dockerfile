FROM golang:1.22-alpine

RUN apk update && apk upgrade && \
    apk --update add git make wget curl jq alpine-sdk

WORKDIR /application

COPY . .

RUN mkdir bin
RUN make build-survey-microservice