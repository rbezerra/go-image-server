FROM golang:1.13-alpine3.10

RUN mkdir /app

ADD . /app

WORKDIR /app

RUN apk update && \
apk add build-base && \
apk add curl git vim wget

RUN go get -d
RUN go build -o main . 

CMD [ "/app/main" ]