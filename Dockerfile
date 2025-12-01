ARG NODE_VERSION=25.2.1
ARG GO_VERSION=1.25.4

FROM node:${NODE_VERSION}-alpine AS node

FROM golang:${GO_VERSION}-alpine AS builder

RUN apk add --no-cache git make

COPY --from=node /usr/local/bin/node /usr/local/bin/node
COPY --from=node /usr/local/bin/npm /usr/local/bin/npm
COPY --from=node /usr/local/bin/npx /usr/local/bin/npx
COPY --from=node /usr/local/lib/node_modules /usr/local/lib/node_modules

WORKDIR /app
COPY . .

RUN make gen
RUN go build -ldflags="-w -s" -o app

FROM alpine:latest

RUN apk add --no-cache ca-certificates

WORKDIR /app
COPY --from=builder /app/app .

ENTRYPOINT ["./app"]
