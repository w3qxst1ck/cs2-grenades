# FROM ubuntu:22.04
# RUN echo 'APT::Install-Suggests "0";' >> /etc/apt/apt.conf.d/00-docker
# RUN echo 'APT::Install-Recommends "0";' >> /etc/apt/apt.conf.d/00-docker

# WORKDIR /build

# COPY ./bin/api .
# COPY .env .
# COPY ./migrations ./migrations

# RUN apt-get update

# CMD ["./api"]

FROM golang:1.20.3

WORKDIR /build

COPY ./bin/api .
COPY .env .
COPY ./migrations ./migrations
COPY entrypoint.sh .

RUN apt-get update
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# CMD ["./api"]
#CMD ["migrate -path=./migrations -database=$GRENADES_DB_DSN up"]