<h1 align='center'>wormholes</h1>
<p align="center">
  <img alt="docker status" src="https://github.com/mohitsinghs/wormholes/actions/workflows/docker.yml/badge.svg">
</p>
<p align="center">
  <b>A self-hosted and scalable link shortener</b><br/>
  <sub>powered by Fiber, PostgreSQL, TimescaleDB and Redis</sub>
</p>
<br />

## Setting Up

Until I figure out Docker and K8s config, you can test this with [wh-scripts](https://github.com/mohitsinghs/wh-scripts).

## Customizing

All three services are highly customizable with environment variables. The default values are there to make them work on my local setup, for production these needs modifications.

### Director

| Name          | Purpose                             |                                                Default |
| ------------- | ----------------------------------- | -----------------------------------------------------: |
| `PORT`        | The port to run                     |                                                 `5000` |
| `BATCH_SIZE`  | Size of batch when ingesting events |                                                `10000` |
| `STREAM`      | Number of streams to ingest events  |                                                    `8` |
| `PG_URI`      | URI for connecting to PostgreSQL    | `postgres://postgres:postgres@localhost:5432/postgres` |
| `PG_MAX_CONN` | Maximum connections for PostgreSQL  |                                                 `5000` |
| `REDIS_URI`   | URI for connecting to Redis         |                      `redis://:redis@localhost:6379/0` |
| `TS_URI`      | URI for connecting to TimescaleDB   | `postgres://postgres:postgres@localhost:5433/postgres` |

### Generator

| Name              | Purpose                                   |                                                Default |
| ----------------- | ----------------------------------------- | -----------------------------------------------------: |
| `PORT`            | The port to run                           |                                                 `5001` |
| `ID_SIZE`         | Size of generated IDs                     |                                                    `7` |
| `MAX_LIMIT`       | Limit of IDs to be stored                 |                                              `1000000` |
| `ERROR_RATE`      | Rate of false positives in bloom filter   |                                            `0.0000001` |
| `BUCKET_SIZE`     | Number of buckets to store IDs            |                                                   `8 ` |
| `BUCKET_CAPACITY` | Number of IDs to store in a single bucket |                                              `100000 ` |
| `PG_URI`          | URI for connecting to PostgreSQL          | `postgres://postgres:postgres@localhost:5432/postgres` |
| `PG_MAX_CONN`     | Maximum connections for PostgreSQL        |                                                 `5000` |

### Creator

| Name          | Purpose                            |                                                Default |
| ------------- | ---------------------------------- | -----------------------------------------------------: |
| `PORT`        | The port to run                    |                                                 `5002` |
| `GEN_ADDR`    | Address of generator instance      |                                       `localhost:5001` |
| `BATCH_SIZE`  | Size of batch when ingesting links |                                                `10000` |
| `PG_URI`      | URI for connecting to PostgreSQL   | `postgres://postgres:postgres@localhost:5432/postgres` |
| `PG_MAX_CONN` | Maximum connections for PostgreSQL |                                                 `5000` |
| `REDIS_URI`   | URI for connecting to Redis        |                      `redis://:redis@localhost:6379/0` |
