#!/bin/env bash

# run postgres
docker run -d \
  --name=postgres \
  -p 127.0.0.1:5432:5432 \
  -e POSTGRES_PASSWORD=postgres \
  -v "$PWD/deploy/conf/postgres.conf":/etc/postgresql/postgresql.conf \
  postgres:alpine \
  -c config_file=/etc/postgresql/postgresql.conf

# run timescale
docker run -d \
  --name=timescale \
  -p 127.0.0.1:5433:5432 \
  -e POSTGRES_PASSWORD=postgres \
  timescale/timescaledb:latest-pg13

# run redis
docker run -d \
  --name=redis \
  -p 127.0.0.1:6379:6379 \
  -v "$PWD/deploy/conf/redis.conf":/usr/local/etc/redis/redis.conf \
  redis:6-alpine \
  redis-server /usr/local/etc/redis/redis.conf
