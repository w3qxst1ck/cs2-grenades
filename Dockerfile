FROM golang:alpine AS builder

WORKDIR /build

ADD go.mod .

COPY . .

RUN GOOS=linux go build -o=./bin/api ./cmd/api


FROM alpine

WORKDIR /build

COPY .env .

COPY --from=builder /build/bin/api /build/bin/api

RUN mkdir -p internal/images

CMD ["./bin/api"]
