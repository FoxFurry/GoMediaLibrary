FROM golang:1.16-alpine

RUN apk add --no-cache make curl gcc libc-dev postgresql-client

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . ./

CMD make integration-tests