# syntax=docker/dockerfile:1
FROM golang:1.21.1-alpine AS build

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

FROM scratch

COPY --from=build /usr/src/app/swagger.yaml /swagger.yaml
COPY --from=build /usr/local/bin/segmentify /segmentify

ENTRYPOINT ["/segmentify"]