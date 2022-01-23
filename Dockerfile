# syntax=docker/dockerfile:1

##
## Build
##
FROM golang:1.17-buster AS build
WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o /api-example

##
## Deploy
##
FROM gcr.io/distroless/base-debian10

WORKDIR /app

COPY --from=build /api-example /app/api-example
## This configurations are only used for local development, the package github.com/spf13/viper override the app.env env vars by the docker container env vars... 
COPY --from=build ./app/app.env /app/

EXPOSE 22210

USER nonroot:nonroot

ENTRYPOINT ["/app/api-example"]