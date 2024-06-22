FROM golang:alpine AS builder

WORKDIR /build

ADD go.mod .

COPY . .

# RUN go get -u -d github.com/golang-migrate/migrate/cmd/migrate
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
# RUN migrate -path=./migrations -database='postgres://admin:w3qxst1ck@postgresdb/grenades?sslmode=disable' up

RUN GOOS=linux go build -o=./bin/api ./cmd/api

CMD ["sh", "./migrations.sh"]

FROM alpine

WORKDIR /build

COPY .env .

COPY entrypoint.sh .

COPY --from=builder /build/bin/api /build/bin/api

RUN mkdir -p internal/images

CMD ["./bin/api"]
