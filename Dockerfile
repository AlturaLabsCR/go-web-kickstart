ARG NODE_VERSION=25.2.1
ARG GO_VERSION=1.25.4
ARG CGO_ENABLED=0

FROM node:${NODE_VERSION}-alpine AS node

FROM golang:${GO_VERSION}-alpine AS builder

RUN apk add --no-cache libstdc++ libgcc make

COPY --from=node /usr/local /usr/local

WORKDIR /app
COPY . .

RUN make gen
RUN CGO_ENABLED="${CGO_ENABLED}" go build -ldflags="-w -s" -o app

FROM alpine:latest

RUN apk add --no-cache ca-certificates

WORKDIR /app
COPY --from=builder /app/app .

ENTRYPOINT ["./app"]
