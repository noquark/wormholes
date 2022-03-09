<h1 align="center">wormholes</h1>
<p align='center'>
  <img alt="GitHub Workflow Status" src="https://img.shields.io/github/workflow/status/mohitsinghs/wormholes/docker?logo=github&style=for-the-badge" />
  <img alt="Go Report Card" src="https://img.shields.io/badge/go%20report-A-green.svg?style=for-the-badge&logo=go&logoColor=white" />
  <img alt="GitHub" src="https://img.shields.io/github/license/mohitsinghs/wormholes?logo=gnu&style=for-the-badge" />
  <img alt="Go Version" src="https://img.shields.io/github/go-mod/go-version/mohitsinghs/wormholes?logo=go&logoColor=white&style=for-the-badge" />
</p>
<br />

## Features

- Fail Safe
- Scalable and Distributed
- Highly Customizable
- High Performance
- API Access

## Why ?

There are a lot of link-shorteners out there, open-source and commercial. they provide easy access to resources in diverse fields and it's an easy problem to solve until it's not.
When I witnessed how a quick in-house implementation can cost a business and will refuse to scale, I was curious on how to scale these reliably. When a quick lookup didn't help me find any solution, I decided to write my own. I've detailed about this [here](https://mohitsingh.in/building-a-link-shortner) and [here](https://mohitsingh.in/a-distributed-link-shortner).

## Getting started

### Run databases

Test config for redis and postgres are inside `deploy/conf`. Included postgres config is tuned for my system for 5000 connections. Generate your own with [pgtune](https://pgtune.leopard.in.ua/#/). Run postgres, timescale and redis now &mdash;

```sh
./deploy/start_db.sh
```

After this, you should have three running in docker containers.

### Run with docker

```sh
docker run -d --network host --name generator ghcr.io/mohitsinghs/wormholes:latest
docker run -d --network host --name director ghcr.io/mohitsinghs/wormholes:latest ./wormholes -as director
docker run -d --network host --name creator ghcr.io/mohitsinghs/wormholes:latest ./wormholes -as creator

```

### Or, Run manually

Finally run generator followed by director and creator

```sh
# clone repository
git clone https://github.com/mohitsinghs/wormholes
cd wormholes
# build binary
go build .
# run all three services
./wormholes -as generator
./wormholes -as creator
./wormholes -as director
```

## Load Testing

Make sure everything is running and accessible. Now, install [wrk](https://github.com/wg/wrk) in your system. We will be using this to load test our setup.

```sh
# load test link creation
wrk -t8 -d10s -c100 -s "./deploy/load/put.lua" http://localhost:5002/api/v1/links
# load test link data API. Get one of shortIDs created in previous step
wrk -t8 -d10s -c100 http://localhost:5001/api/v1/links/<shortID>
# load test link redirection
wrk -t8 -d10s -c100 -s "./deploy/load/put.lua"  http://localhost:5000/<shortID>
```

## Environment Variables

Wormholes are highly customizable with environment variables. The default values are there to make them work on my local setup, for production these needs modifications. Following are the environment variables and there usage &mdash;

| Name          | Purpose                                   |                                           Default Value |
| ------------- | ----------------------------------------- | ------------------------------------------------------: |
| `PORT`        | The port to run                           | `5000` (director), `5001` (generator), `5002` (creator) |
| `BATCH_SIZE`  | Size of batch when ingesting events       |                                                 `10000` |
| `STREAMS`     | Number of streams to ingest events        |                                                     `8` |
| `ID_SIZE`     | Size of generated IDs                     |                                                     `7` |
| `BLOOM_MAX`   | Limit of IDs to be stored                 |                                               `1000000` |
| `BLOOM_ERROR` | Rate of false positives in bloom filter   |                                             `0.0000001` |
| `BUCKET_SIZE` | Number of buckets to store IDs            |                                                     `8` |
| `BUCKET_CAP`  | Number of IDs to store in a single bucket |                                               `100000 ` |
| `GEN_ADDR`    | Address of generator instance             |                                        `localhost:5001` |
| `TS_URI`      | URI for connecting to TimescaleDB         |  `postgres://postgres:postgres@localhost:5433/postgres` |
| `PG_URI`      | URI for connecting to PostgreSQL          |  `postgres://postgres:postgres@localhost:5432/postgres` |
| `PG_MAX_CONN` | Maximum connections for PostgreSQL        |                                                  `5000` |
| `REDIS_URI`   | URI for connecting to Redis               |                       `redis://:redis@localhost:6379/0` |

## Behind the name

Wormholes are tunnels made by earthworms and a link between two disparate points in spacetime. Short links are essentially a links to digital thing. This project is inspired by both and hence named after them.

## Having problems ?

Feel free to open an issue or feature request. I'll try to help you as long as you are being reasonable.
