#!/bin/sh
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
migrate -path=./migrations -database='postgres://admin:w3qxst1ck@postgresdb/grenades?sslmode=disable' up