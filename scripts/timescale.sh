#!/bin/env zsh
# Start timescaledb and bind to local port 5433

docker run -d \
  --name=timescale \
  -p 5433:5432 \
  -e POSTGRES_PASSWORD=postgres \
  timescale/timescaledb:latest-pg13
