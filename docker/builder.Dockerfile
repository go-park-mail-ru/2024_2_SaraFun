FROM golang:1.22-alpine

RUN apk update && apk upgrade && \
    apk --update add git make wget curl jq alpine-sdk

WORKDIR /application

COPY . .

RUN mkdir bin
RUN make build-sparkit
RUN make build-auth-microservice
RUN make build-personalities-microservice
RUN make build-communications-microservice
RUN make build-message-microservice
RUN make build-survey-microservice
RUN make build-payments-microservice
