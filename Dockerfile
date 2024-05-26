FROM golang:alpine AS builder

WORKDIR /build

ADD go.mod .

COPY . .

RUN GOOS=linux go build -o=./bin/api ./cmd/api

FROM alpine

WORKDIR /build

COPY .env .

COPY entrypoint.sh .

COPY --from=builder /build/bin/api /build/bin/api
