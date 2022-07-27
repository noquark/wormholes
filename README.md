<h1 align="center">wormholes</h1>
<p align='center'>
  <img alt="GitHub" src="https://img.shields.io/github/workflow/status/mohitsinghs/wormholes/docker?logo=github&style=flat-square" />
  <img alt="Go Report Card" src="https://img.shields.io/badge/go%20report-A-green.svg?style=flat-square&logo=go&logoColor=white" />
  <img alt="GitHub" src="https://img.shields.io/github/license/mohitsinghs/wormholes?logo=gnu&style=flat-square" />
  <img alt="Go Version" src="https://img.shields.io/github/go-mod/go-version/mohitsinghs/wormholes?logo=go&logoColor=white&style=flat-square" />
</p>
<br />

## Why wormholes ?

I was curious on how to scale link-shortners reliably and decided to write one. See [Building a link shortner](https://mohitsingh.in/code/building-a-link-shortner) for the story.

## Getting Started

1. Default database configs are inside `deploy/conf`. Verify those.
2. Included postgres config is tuned for my system for 5000 connections. Generate your own with [pgtune](https://pgtune.leopard.in.ua/#/).
3. Run following to start postgres, timescale and redis in containers &mdash;

```sh
./deploy/start_db.sh
```

4. Run wormholes with docker

```sh
docker run -d --network host --name generator ghcr.io/mohitsinghs/wormholes:latest
```

Or, build and run the binary yourself

```sh
# clone repository
git clone https://github.com/mohitsinghs/wormholes
cd wormholes
# build binary
go build .
# run unified
./wormholes
```

## Using in production

You can start each of generator, creator and director independently using same image with different `-as` flag. This allows you to scale creator and director as per your needs.

## Environment Variables

Wormholes is highly customizable with environment variables. The default values are there to make them work on local setup, for production these needs modifications.

Following are the environment variables and there usage &mdash;

1. For distributed setup

| Name          | Purpose                                   |                                           Default Value |
| ------------- | ----------------------------------------- | ------------------------------------------------------: |
| `PORT`        | The port to run (`5000` when unified)     | `5000` (director), `5001` (generator), `5002` (creator) |
| `BATCH_SIZE`  | Size of batch when ingesting events       |                                                 `10000` |
| `STREAMS`     | Number of streams to ingest events        |                                                     `8` |
| `ID_SIZE`     | Size of generated IDs                     |                                                     `7` |
| `BLOOM_MAX`   | Limit of IDs to store                     |                                               `1000000` |
| `BLOOM_ERROR` | Rate of false positives in bloom filter   |                                             `0.0000001` |
| `BUCKET_SIZE` | Number of buckets to store IDs            |                                                     `8` |
| `BUCKET_CAP`  | Number of IDs to store in a single bucket |                                               `100000 ` |
| `GEN_ADDR`    | Address of generator instance             |                                        `localhost:5001` |
| `TS_URI`      | URI for connecting to TimescaleDB         |  `postgres://postgres:postgres@localhost:5433/postgres` |
| `PG_URI`      | URI for connecting to PostgreSQL          |  `postgres://postgres:postgres@localhost:5432/postgres` |
| `PG_MAX_CONN` | Max connections for PostgreSQL            |                                                  `5000` |
| `REDIS_URI`   | URI for connecting to Redis               |                       `redis://:redis@localhost:6379/0` |

2. For unified setup

| Name          | Purpose                                   |             Default Value |
| ------------- | ----------------------------------------- | ------------------------: |
| `PORT`        | The port to run                           |                    `5000` |
| `BATCH_SIZE`  | Size of batch when ingesting events       |                   `10000` |
| `STREAMS`     | Number of streams to ingest events        |                       `8` |
| `ID_SIZE`     | Size of generated IDs                     |                       `7` |
| `BLOOM_MAX`   | Limit of IDs to store                     |                 `1000000` |
| `BLOOM_ERROR` | Rate of false positives in bloom filter   |               `0.0000001` |
| `BUCKET_SIZE` | Number of buckets to store IDs            |                       `8` |
| `BUCKET_CAP`  | Number of IDs to store in a single bucket |                 `100000 ` |
| `TS_URI`      | URI for connecting to TimescaleDB         | same as distributed setup |
| `PG_URI`      | URI for connecting to PostgreSQL          | same as distributed setup |
| `PG_MAX_CONN` | Max connections for PostgreSQL            |                    `5000` |
| `REDIS_URI`   | URI for connecting to Redis               | same as distributed setup |

## Load Testing with wrk

### Requirements

1. Everything is running in distributed mode.
2. [wrk](https://github.com/wg/wrk) is installed in your system.

### Tests

1. Load test link creation

```sh
wrk -t8 -d10s -c100 -s "./deploy/load/put.lua" http://localhost:5002/api/v1/links
```

2.  Load test link data API. Get one of shortIDs created in previous step

```sh
wrk -t8 -d10s -c100 http://localhost:5001/api/v1/links/<shortID>
```

3. load test link redirection

```sh
wrk -t8 -d10s -c100 http://localhost:5000/<shortID>
```

## Contributing

Feel free to open an issue or pull request.
