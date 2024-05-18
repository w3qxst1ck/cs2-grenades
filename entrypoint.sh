#!/bin/bash
migrate -path=./migrations -database=$GRENADES_DB_DSN up
./api