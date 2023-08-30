# syntax=docker/dockerfile:1
FROM golang:1.21.0-alpine

WORKDIR /usr/src/app

#install go-swagger
RUN go install github.com/go-swagger/go-swagger/cmd/swagger@latest

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /usr/local/bin/segmentify

# generate swagger file
RUN swagger generate spec -o ./swagger.yaml --scan-models

# Set environment variables
ENV DB_CONNECTION_STRING user=postgres dbname=postgres host=192.168.0.1 port=5432 sslmode=disable

CMD ["segmentify"]